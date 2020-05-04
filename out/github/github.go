package github

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Repo struct {
	Repository string
	Location string
	AccessToken string
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client
type Client interface {
	ExecuteGithubApi(string, string, string, []byte) ([]byte, error)
	ExecuteGithubCmd(...string) (string, error)
}

type GithubClient struct {}

func (g GithubClient) ExecuteGithubApi(url string, method string, authorizationHeaders string, body []byte) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token " + authorizationHeaders)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git api command %w", err)
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read from git api command %w: %s", err, response)
	}

	return response, nil
}


func (g GithubClient) ExecuteGithubCmd(param ...string) (string, error){

	output, err := exec.Command("git", param...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execute git command %w: %s", err, string(output))
	}

	return string(output), nil
}