package pullRequest

import (
	"encoding/json"
	"fmt"
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
	var fakeClient *fakes.FakeClient

	it.Before(func() {
		fakeClient = &fakes.FakeClient{}
	})

	when("auto merge is false", func() {
		it("make a branch and push to remote", func() {
			description := `this is the description of the PR
	it may have new lines and
	differennt @ or # kind of characters
	it might also have / or \ `
			pr := NewPullRequest(description, "A new PR", "master", "pr-by-test", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			repo := "test-org/test"
			location := "artifacts"
			accessToken := "123456789"
			branchName, prNumber, err := pr.CreatePullRequest(repo, "", location, accessToken, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, _, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo))
			assert.Equal(t, header, accessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", location, "push", "https://123456789:x-oauth-basic@github.com/test-org/test.git", "--no-verify"})

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, jsonActualBody["title"], "A new PR")
			assert.Equal(t, jsonActualBody["body"], description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], "test-org:"+branchName)
		})
		it("apply the default values if not provided", func() {
			pr := NewPullRequest("", "", "", "", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			repo := "test-org/test"
			location := "artifacts"
			accessToken := "123456789"
			branchName, prNumber, err := pr.CreatePullRequest(repo, "", location, accessToken, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, _, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo))
			assert.Equal(t, header, accessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", location, "push", "https://123456789:x-oauth-basic@github.com/test-org/test.git", "--no-verify"})

			assert.Equal(t, jsonActualBody["title"], "Pull request by bot")
			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, jsonActualBody["body"], "This is default description of the PR")
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], "test-org:"+branchName)
		})
		it("apply the few default values if not provided", func() {
			pr := NewPullRequest("This is the description of new PR by test", "Changing the version string to 12 from 14", "", "", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			repo := "test-org/test"
			location := "artifacts"
			accessToken := "123456789"
			branchName, prNumber, err := pr.CreatePullRequest(repo, "", location, accessToken, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, method, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo))
			assert.Equal(t, header, accessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", location, "push", "https://123456789:x-oauth-basic@github.com/test-org/test.git", "--no-verify"})

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, method, "POST")
			assert.Equal(t, jsonActualBody["title"], pr.Title)
			assert.Equal(t, jsonActualBody["body"], pr.Description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], "test-org:"+branchName)
		})
	})
	when("auto merge is true", func() {
		it("make a branch and push to remote", func() {
			description := `this is the description of the PR
	it may have new lines and
	differennt @ or # kind of characters
	it might also have / or \ `
			pr := NewPullRequest(description, "A new PR", "master", "pr-by-test", true)

			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			repo := "test-org/test"
			location := "artifacts"
			accessToken := "123456789"
			_, prNumber, err := pr.CreatePullRequest(repo, "", location, accessToken, fakeClient)
			require.NoError(t, err)

			url, method, header, body := fakeClient.ExecuteGithubApiArgsForCall(1)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/merge", repo, "12345"))
			assert.Equal(t, header, accessToken)

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, method, "PUT")
			assert.Equal(t, jsonActualBody["commit_title"], fmt.Sprintf("Auto merge pull request 12345"))
			assert.Equal(t, jsonActualBody["commit_message"], "")
			assert.Equal(t, jsonActualBody["sha"], "somesha")
		})
	})
	when("the forked repository is not empty", func() {
		it("pushes the branch to the remote repository and makes the PR from there", func() {
			remoteRepo := "test-org/test"
			forkedRepo := "fork-user/test"
			location := "some-location"
			accessToken := "some-access-token"
			description := `this is the description of the PR
it may have new lines and 
differennt @ or # kind of characters
it might also have / or \ `
			pr := NewPullRequest(description, "A new PR", "master", "pr-by-test", false)
			fakeClient.ExecuteGithubApiReturnsOnCall(0, []byte(`{"number":12345,"head":{"sha":"somesha"}}`), nil)

			branchName, prNumber, err := pr.CreatePullRequest(remoteRepo, forkedRepo, location, accessToken, fakeClient)
			require.NoError(t, err)

			assert.Regexp(t, "^pr-by-test-.*$", branchName)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, _, header, body := fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", remoteRepo))
			assert.Equal(t, header, accessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", location, "checkout", "-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", location, "push", "https://some-access-token:x-oauth-basic@github.com/fork-user/test.git", "--no-verify"})

			assert.Equal(t, prNumber, 12345)
			assert.Equal(t, jsonActualBody["title"], "A new PR")
			assert.Equal(t, jsonActualBody["body"], description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], "fork-user:"+branchName)
		})
	})
	when("pushing the branch fails", func() {
		it("does not show the token in the error message", func() {
			accessToken := "some-access-token"

			fakeClient.ExecuteGithubCmdReturnsOnCall(1, "", fmt.Errorf("error which includes token %s", accessToken))

			pr := NewPullRequest("", "", "", "", false)
			_, _, err := pr.CreatePullRequest("", "", "", accessToken, fakeClient)

			require.Error(t, err)
			assert.NotContains(t, err.Error(), accessToken)
		})
	})
}
