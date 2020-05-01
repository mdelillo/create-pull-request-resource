package pullRequest

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/out/github"
	"time"
)

type PullRequest struct {
	Description  string
	Title  string
	Base         string
	BranchPrefix string
	AutoMerge    bool
}

type RequestBody struct{
	Title		string
	Body		string
	Head		string
	Base		string
}

func (p PullRequest) CreatePullRequest(repo github.Repo, client github.Client) (string, error) {
	branchName, err := p.createRemoteBranch(repo, client)
	if err != nil {
		return branchName, fmt.Errorf("failed to create a reamote branch with name %s %w", branchName, err)
	}
	err = p.createPullRequestFor(branchName, repo, client)
	if err != nil {
		return branchName, fmt.Errorf("failed to create a pull request %w", err)
	}
	return branchName, nil
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

func (p PullRequest) createPullRequestFor(branchName string, repo github.Repo, client github.Client) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository)

	body, _ := json.Marshal(map[string]string{
		"title": p.Title,
		"body":  p.Description,
		"head":  branchName,
		"base":  p.Base,
	})
	apiOutput, err := client.ExecuteGithubApi(url, repo.AccessToken, body)
	if err != nil {
		return fmt.Errorf("failed to POST a pull request %w %s", err, apiOutput)
	}
	return nil
}

