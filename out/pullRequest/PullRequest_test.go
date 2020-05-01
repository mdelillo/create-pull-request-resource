package pullRequest

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal/create-pull-request-resource/out/github"
	fakes "github.com/pivotal/create-pull-request-resource/out/github/githubfakes"
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
			pr := PullRequest{Description: description, Base: "master", BranchPrefix: "my-new-pr", Title: "A new PR", AutoMerge: true}

			branchName, err := pr.CreatePullRequest(repo, fakeClient)
			require.NoError(t, err)

			paramForFirstCmd := fakeClient.ExecuteGithubCmdArgsForCall(0)
			paramForSecondCmd := fakeClient.ExecuteGithubCmdArgsForCall(1)
			url, header , body:= fakeClient.ExecuteGithubApiArgsForCall(0)
			var jsonActualBody map[string]string
			json.Unmarshal(body, &jsonActualBody)

			assert.Equal(t, url, fmt.Sprintf("https://api.github.com/repos/%s/pulls", repo.Repository))
			assert.Equal(t, header, repo.AccessToken)

			assert.EqualValues(t, paramForFirstCmd, []string{"-C", repo.Location, "checkout","-b", branchName})
			assert.EqualValues(t, paramForSecondCmd, []string{"-C", repo.Location, "push", "origin", branchName})

			assert.Equal(t, jsonActualBody["title"], "A new PR")
			assert.Equal(t, jsonActualBody["body"], description)
			assert.Equal(t, jsonActualBody["base"], "master")
			assert.Equal(t, jsonActualBody["head"], branchName)
		})
	})
}
