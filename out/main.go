package main

import (
	"encoding/json"
	"github.com/pivotal/create-pull-request-resource/out/github"
	"github.com/pivotal/create-pull-request-resource/out/pullRequest"
	"log"
	"os"
)

type Source struct {
	GithubToken       	string `json:"github_token"`
	Repository   		string `json:"repository"`
	Base   				string `json:"base,omitempty"`
	Description   		string `json:"description,omitempty"`
	Title		   		string `json:"title,omitempty"`
	BranchPrefix   	    string `json:"branch_prefix,omitempty"`
	AutoMerge   		bool `json:"auto_merge,omitempty"`
}


type OutRequest struct {
	Source Source    `json:"source"`
}

func main() {
	if len(os.Args) != 2 {
		log.Println("usage:", os.Args[0],  "<sources directory>")
		os.Exit(1)
	}

	var request OutRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("failed to read request: %s", err.Error())
	}
	sourcePath := os.Args[1]

	log.Println("starting to create PR ")
	log.Println(request)
	log.Println(os.Args[0])
	log.Println(os.Args[1])

	repo := github.Repo{AccessToken: request.Source.GithubToken, Repository: request.Source.Repository, Location: sourcePath}

	newPullRequest := pullRequest.NewPullRequest(request.Source.Description, request.Source.Title, request.Source.Base, request.Source.BranchPrefix, request.Source.AutoMerge)

	branchName, prNumber, err := newPullRequest.CreatePullRequest(repo, github.GithubClient{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	log.Println("created a PR with Branch name:", branchName, "and pull request number", prNumber)
	if newPullRequest.AutoMerge {
		log.Println("Merged the Pull Request", prNumber)
	}
}
