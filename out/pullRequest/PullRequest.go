package pullRequest

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/out/github"
	"time"
)

type PullRequest struct {
	Description  string
	Title 	 	 string
	Base         string
	BranchPrefix string
	AutoMerge    bool
}

type PrRespeonse struct {
	Number int  `json:"number"`
	Head   Head `json:"head"`
}
type Head struct{
	SHA string `json:"sha"`
}

func NewPullRequest(descrption string, title string,base string,branchPrefix string,autoMerge bool) PullRequest {

	request := PullRequest{Description:descrption, Title:title, Base:base, BranchPrefix:branchPrefix, AutoMerge:autoMerge}

	if descrption == "" {request.Description = "This is default description of the PR"}
	if title == "" {request.Title = "Pull request by bot"}
	if base == "" {request.Base = "master"}
	if branchPrefix == "" {request.BranchPrefix = "pr-by-bot"}

	return request
}

func (p PullRequest) CreatePullRequest(repo github.Repo, client github.Client) (string, int, error) {
	branchName, err := p.createRemoteBranch(repo, client)
	if err != nil {
		return branchName, 0, fmt.Errorf("failed to create a reamote branch with name %s %w", branchName, err)
	}
	prResponse, err := p.createPullRequestFor(branchName, repo, client)
	if err != nil {
		return branchName, 0, fmt.Errorf("failed to create a pull request %w", err)
	}

	if p.AutoMerge {
		output, err := p.mergePullRequest(prResponse, repo, client)
		if err != nil {
			return branchName, 0, fmt.Errorf("failed to merge a pull request %w %s", err, output)
		}
	}
	return branchName, prResponse.Number, nil
}

func (p PullRequest) createRemoteBranch(repo github.Repo, client github.Client) (string, error) {

	branchName := fmt.Sprintf("%s-%d", p.BranchPrefix, time.Now().Unix())

	output, err := client.ExecuteGithubCmd( "-C", repo.Location, "checkout","-b", branchName)
	if err != nil {
		return branchName, fmt.Errorf("failed to checkout new branch %w: %s", err, output)
	}

	output, err = client.ExecuteGithubCmd("-C", repo.Location, "push", "origin" , branchName)
	if err != nil {
		return branchName, fmt.Errorf("failed to push to new branch %w: %s", err, output)
	}
	return branchName, nil
}

func (p PullRequest) createPullRequestFor(branchName string, repo github.Repo, client github.Client) (PrRespeonse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository)

	body, _ := json.Marshal(map[string]string{
		"title": p.Title,
		"body":  p.Description,
		"head":  branchName,
		"base":  p.Base,
	})
	apiOutput, err := client.ExecuteGithubApi(url, "POST",  repo.AccessToken, body)
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

func (p PullRequest) mergePullRequest(prResponse PrRespeonse, repo github.Repo, client github.Client) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/merge", repo.Repository, prResponse.Number)

	body, _ := json.Marshal(map[string]string{
		"commit_title": fmt.Sprintf("Auto merge pull request %d", prResponse.Number),
		"commit_message": "",
		"sha": prResponse.Head.SHA,
	})
	apiOutput, err := client.ExecuteGithubApi(url, "PUT",  repo.AccessToken, body)
	if err != nil {
		return nil, fmt.Errorf("failed to PULL merge for pull request ID %d: %w \n %s", prResponse.Number, err, apiOutput)
	}
	return apiOutput, nil
}

