package main

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/out/github"
	"github.com/pivotal/create-pull-request-resource/out/pullRequest"
	"log"
	"os"
)

type Source struct {
	GithubToken       	string `json:"github_token"`
	Repository   		string `json:"repository"`
	Base   				string `json:"base"`
	Description   		string `json:"description"`
	Title		   		string `json:"title"`
	BranchPrefix   	    string `json:"branch_prefix"`
	Location   			string `json:"location"`
	AutoMerge   		bool `json:"auto_merge"`
}


type OutRequest struct {
	Source Source    `json:"source"`
}

func main() {
	fmt.Printf("starting to create PR \n")
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <sources directory>", os.Args[0])
		os.Exit(1)
	}

	var request OutRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("failed to read request: %s", err.Error())
	}

	repo := github.Repo{AccessToken: request.Source.GithubToken, Repository: request.Source.Repository, Location: request.Source.Location}

	newPullRequest := pullRequest.PullRequest{Description: request.Source.Description, Title:request.Source.Title, BranchPrefix: request.Source.BranchPrefix, Base: request.Source.Base, AutoMerge: request.Source.AutoMerge}

	branchName, err := newPullRequest.CreatePullRequest(repo, github.GithubClient{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.Println("craeted a PR with Branch name:", branchName)
}
