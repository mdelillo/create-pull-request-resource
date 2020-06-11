package pullRequest

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal/create-pull-request-resource"
	"github.com/pivotal/create-pull-request-resource/github"
	"strings"
	"time"
)

type PullRequest struct {
	Description  string
	Title        string
	Base         string
	BranchPrefix string
	AutoMerge    bool
}

func NewPullRequest(description string, title string, base string, branchPrefix string, autoMerge bool) PullRequest {

	request := PullRequest{Description: description, Title: title, Base: base, BranchPrefix: branchPrefix, AutoMerge: autoMerge}

	if description == "" {
		request.Description = "This is default description of the PR"
	}
	if title == "" {
		request.Title = "Pull request by bot"
	}
	if base == "" {
		request.Base = "master"
	}
	if branchPrefix == "" {
		request.BranchPrefix = "pr-by-bot"
	}

	return request
}

func (p PullRequest) CreatePullRequest(remoteRepo, forkedRepo, location, accessToken string, client github.Client) (string, int, error) {
	sourceRepo := remoteRepo
	if forkedRepo != "" {
		sourceRepo = forkedRepo
	}

	branchName, err := p.createRemoteBranch(sourceRepo, location, accessToken, client)
	if err != nil {
		return branchName, 0, fmt.Errorf("failed to create a reamote branch with name %s %w", branchName, err)
	}

	prResponse, err := p.createPullRequestFor(branchName, sourceRepo, remoteRepo, accessToken, client)
	if err != nil {
		return branchName, 0, fmt.Errorf("failed to create a pull request %w", err)
	}

	if p.AutoMerge {
		output, err := p.mergePullRequest(prResponse, remoteRepo, accessToken, client)
		if err != nil {
			return branchName, 0, fmt.Errorf("failed to merge a pull request %w %s", err, output)
		}
	}
	return branchName, prResponse.Number, nil
}

func (p PullRequest) createRemoteBranch(repo, location, accessToken string, client github.Client) (string, error) {
	branchName := fmt.Sprintf("%s-%d", p.BranchPrefix, time.Now().Unix())

	output, err := client.ExecuteGithubCmd("-C", location, "checkout", "-b", branchName)
	if err != nil {
		return branchName, fmt.Errorf("failed to checkout new branch %w: %s", err, output)
	}

	output, err = client.ExecuteGithubCmd("-C", location, "push", fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s.git", accessToken, repo), "--no-verify")
	if err != nil {
		return branchName, fmt.Errorf("failed to push to new branch %w: %s", err, output)
	}
	return branchName, nil
}

func (p PullRequest) createPullRequestFor(branchName string, sourceRepo, remoteRepo, accessToken string, client github.Client) (PrRespeonse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", remoteRepo)

	githubUser := strings.Split(sourceRepo, "/")[0]
	branchName = fmt.Sprintf("%s:%s", githubUser, branchName)
	body, _ := json.Marshal(map[string]string{
		"title": p.Title,
		"body":  p.Description,
		"head":  branchName,
		"base":  p.Base,
	})
	apiOutput, err := client.ExecuteGithubApi(url, "POST", accessToken, body)
	if err != nil {
		return PrRespeonse{}, fmt.Errorf("failed to POST a pull request %w %s", err, apiOutput)
	}

	var pullRequestContent PrRespeonse
	err = json.Unmarshal(apiOutput, &pullRequestContent)
	if err != nil {
		return PrRespeonse{}, fmt.Errorf("failed to unmarshall create pull request's response: %s", err)
	}

	return pullRequestContent, nil
}

func (p PullRequest) mergePullRequest(prResponse PrRespeonse, repo, accessToken string, client github.Client) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/merge", repo, prResponse.Number)

	body, _ := json.Marshal(map[string]string{
		"commit_title":   fmt.Sprintf("Auto merge pull request %d", prResponse.Number),
		"commit_message": "",
		"sha":            prResponse.Head.SHA,
	})
	apiOutput, err := client.ExecuteGithubApi(url, "PUT", accessToken, body)
	if err != nil {
		return nil, fmt.Errorf("failed to PULL merge for pull request ID %d: %w \n %s", prResponse.Number, err, apiOutput)
	}
	return apiOutput, nil
}
