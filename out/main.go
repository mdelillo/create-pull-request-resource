package main

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/out/github"
	"github.com/pivotal/create-pull-request-resource/out/pullRequest"
	"log"
	"os"
	"path/filepath"
)

type Source struct {
	GithubToken      string `json:"github_token"`
	RemoteRepository string `json:"remote_repository"`
}
type Params struct {
	RepoLocation 	 string `json:"repo_location"`
	Base             string `json:"base,omitempty"`
	Description      string `json:"description,omitempty"`
	Title            string `json:"title,omitempty"`
	BranchPrefix     string `json:"branch_prefix,omitempty"`
	AutoMerge        bool   `json:"auto_merge,omitempty"`
}


type OutRequest struct {
	Source Source    `json:"source"`
	Params Params    `json:"params"`
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

	repo := github.Repo{AccessToken: request.Source.GithubToken, Repository: request.Source.RemoteRepository, Location: filepath.Join(os.Args[1], request.Params.RepoLocation)}

	newPullRequest := pullRequest.NewPullRequest(request.Params.Description, request.Params.Title, request.Params.Base, request.Params.BranchPrefix, request.Params.AutoMerge)

	branchName, prNumber, err := newPullRequest.CreatePullRequest(repo, github.GithubClient{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	outPutResponse := fmt.Sprintf(`{"breanchName":"%s","prNumber":%d,"merged":"%t"}`, branchName, prNumber, newPullRequest.AutoMerge)
	fmt.Println(string(outPutResponse))
}
