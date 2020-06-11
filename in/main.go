package main

import (
	"encoding/json"
	"fmt"
	. "github.com/pivotal/create-pull-request-resource"
	"github.com/pivotal/create-pull-request-resource/github"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <sources directory>", os.Args[0])
		os.Exit(1)
	}

	var request InRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		log.Fatalf("failed to read request: %s", err.Error())
	}

	client := github.GithubClient{}

	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", request.Source.RemoteRepository, request.Version.Ref)

	apiOutput, err := client.ExecuteGithubGetApi(url, request.Source.GithubToken)
	if err != nil {
		log.Fatalf("Could not make a request for listing the newly create PR: %s", err.Error())
	}

	var pullRequestContent PrRespeonse

	err = json.Unmarshal(apiOutput, &pullRequestContent)
	if err != nil {
		log.Fatalf("failed to unmarshall get pull request's response: %s", err)
	}

	inPutResponse := fmt.Sprintf(`{ "version": { "ref": "%s" },"metadata": [{"sha":"%s"}]}`, request.Version.Ref, pullRequestContent.Head.SHA)

	fmt.Println(string(inPutResponse))
}
