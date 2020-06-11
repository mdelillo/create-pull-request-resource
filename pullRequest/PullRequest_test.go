package pullRequest

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/github"
	fakes "github.com/pivotal/create-pull-request-resource/github/githubfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePullRequestTask(t *testing.T) {
	spec.Run(t, "CreatePullRequest", testCreatePullRequestTask, spec.Report(report.Terminal{}))
}

func testCreatePullRequestTask(t *testing.T, when spec.G, it spec.S) {
	when("auto merge is false", func() {
		it("make a branch and push to remote", func() {
			var fakeClient *fakes.FakeClient
			fakeClient = &fakes.FakeClient{}

			repo := github.Repo{"test/test", "artifacts", "123456789"}
			description := `this is the description of the PR
it may have new lines and 
differennt @ or # kind of characters
it might also have / or \ `
			pr := NewPullRequest(description, "A new PR", "master", "pr-by-test", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			branchName, prNumber, err := pr.CreatePullRequest(repo, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, _, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository))
			assert.Equal(t, header, repo.AccessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", repo.Location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", repo.Location, "push", "https://123456789:x-oauth-basic@github.com/test/test.git", "--no-verify"})

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, jsonActualBody["title"], "A new PR")
			assert.Equal(t, jsonActualBody["body"], description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], branchName)
		})
		it("apply the default values if not provided", func() {
			var fakeClient *fakes.FakeClient
			fakeClient = &fakes.FakeClient{}

			repo := github.Repo{"test/test", "artifacts", "123456789"}
			pr := NewPullRequest("", "", "", "", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			branchName, prNumber, err := pr.CreatePullRequest(repo, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, _, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository))
			assert.Equal(t, header, repo.AccessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", repo.Location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", repo.Location, "push", "https://123456789:x-oauth-basic@github.com/test/test.git", "--no-verify"})

			assert.Equal(t, jsonActualBody["title"], "Pull request by bot")
			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, jsonActualBody["body"], "This is default description of the PR")
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], branchName)
		})
		it("apply the few default values if not provided", func() {
			var fakeClient *fakes.FakeClient
			fakeClient = &fakes.FakeClient{}

			repo := github.Repo{"test/test", "artifacts", "123456789"}
			pr := NewPullRequest("This is the description of new PR by test", "Changing the version string to 12 from 14", "", "", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			branchName, prNumber, err := pr.CreatePullRequest(repo, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, method, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository))
			assert.Equal(t, header, repo.AccessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", repo.Location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", repo.Location, "push", "https://123456789:x-oauth-basic@github.com/test/test.git", "--no-verify"})

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, method, "POST")
			assert.Equal(t, jsonActualBody["title"], pr.Title)
			assert.Equal(t, jsonActualBody["body"], pr.Description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], branchName)
		})
	})
	when("auto merge is true", func() {
		it("make a branch and push to remote", func() {
			var fakeClient *fakes.FakeClient
			fakeClient = &fakes.FakeClient{}

			repo := github.Repo{"test/test", "artifacts", "123456789"}
			description := `this is the description of the PR
it may have new lines and 
differennt @ or # kind of characters
it might also have / or \ `
			pr := NewPullRequest(description, "A new PR", "master", "pr-by-test", true)

			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			_, prNumber, err := pr.CreatePullRequest(repo, fakeClient)
			require.NoError(t, err)

			url, method, header, body := fakeClient.ExecuteGithubApiArgsForCall(1)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/merge", repo.Repository, "12345"))
			assert.Equal(t, header, repo.AccessToken)

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, method, "PUT")
			assert.Equal(t, jsonActualBody["commit_title"], fmt.Sprintf("Auto merge pull request 12345"))
			assert.Equal(t, jsonActualBody["commit_message"], "")
			assert.Equal(t, jsonActualBody["sha"], "somesha")
		})
	})
}
