package github

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client
type Client interface {
	ExecuteGithubApi(string, string, string, []byte) ([]byte, error)
	ExecuteGithubCmd(...string) (string, error)
}

type GithubClient struct{}

func (g GithubClient) ExecuteGithubApi(url string, method string, authorizationHeaders string, body []byte) ([]byte, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+authorizationHeaders)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git api command %s", removeSecretsFromOutputs(err.Error(), authorizationHeaders))
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read from git api command %s: %s", removeSecretsFromOutputs(err.Error(), authorizationHeaders), removeSecretsFromOutputs(string(response), authorizationHeaders))
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("error: status code: %d, response: %s", resp.StatusCode, response)
	}

	return response, nil
}

func removeSecretsFromOutputs(content string, secret string) string {
	return strings.Replace(content, secret, "<auth-token>", -1)
}

func (g GithubClient) ExecuteGithubCmd(param ...string) (string, error) {

	output, err := exec.Command("git", param...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execute git command %w: %s", err, string(output))
	}

	return string(output), nil
}

func (g GithubClient) ExecuteGithubGetApi(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute git get call %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read from git api call %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("error: status code: %d, response: %s", resp.StatusCode, body)
	}

	return body, nil
}
