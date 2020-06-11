package main

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal/create-pull-request-resource"
	"github.com/pivotal/create-pull-request-resource/github"
	"github.com/pivotal/create-pull-request-resource/pullRequest"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("usage:", os.Args[0], "<sources directory>")
		os.Exit(1)
	}

	var request OutRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("failed to read request: %s", err.Error())
	}

	remoteRepo := github.Repo{AccessToken: request.Source.GithubToken, Repository: request.Source.RemoteRepository, Location: filepath.Join(os.Args[1], request.Params.RepoLocation)}
	forkedRepo := github.Repo{AccessToken: request.Source.GithubToken, Repository: request.Source.ForkedRepository, Location: filepath.Join(os.Args[1], request.Params.RepoLocation)}

	newPullRequest := pullRequest.NewPullRequest(request.Params.Description, request.Params.Title, request.Params.Base, request.Params.BranchPrefix, request.Params.AutoMerge)

	branchName, prNumber, err := newPullRequest.CreatePullRequest(remoteRepo, forkedRepo, github.GithubClient{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	outPutResponse := fmt.Sprintf(`{ "version": { "ref": "%d" },"metadata": [{"breanchName":"%s","prNumber":%d,"merged":"%t"}]}`, prNumber, branchName, prNumber, newPullRequest.AutoMerge)
	fmt.Println(string(outPutResponse))
}
