package github_test

import (
	"fmt"
	"github.com/pivotal/create-pull-request-resource/github"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGithub(t *testing.T) {
	spec.Run(t, "TestGithub", testGithub, spec.Report(report.Terminal{}))
}

func testGithub(t *testing.T, when spec.G, it spec.S) {
	var (
		client       github.GithubClient
		githubServer *httptest.Server
	)

	it.Before(func() {
		client = github.GithubClient{}

		githubServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/200":
				w.WriteHeader(200)
				fmt.Fprint(w, "response from github")
			case "/400":
				w.WriteHeader(400)
				fmt.Fprint(w, "error message from github")
			default:
				require.Fail(t, "Unexpected path: %s", r.URL.Path)
			}
		}))
	})

	it.After(func() {
		githubServer.Close()
	})

	when("ExecuteGithubApi", func() {
		when("github returns a 2XX response", func() {
			it("returns the response body", func() {
				response, err := client.ExecuteGithubApi(githubServer.URL+"/200", "GET", "", nil)

				require.NoError(t, err)
				assert.Equal(t, "response from github", string(response))
			})
		})

		when("github returns a non-2XX response", func() {
			it("returns an error", func() {
				_, err := client.ExecuteGithubApi(githubServer.URL+"/400", "GET", "", nil)

				assert.Error(t, err)
				assert.Equal(t, "error: status code: 400, response: error message from github", err.Error())
			})
		})
	})

	when("ExecuteGithubGetApi", func() {
		when("github returns a 2XX response", func() {
			it("returns the response body", func() {
				response, err := client.ExecuteGithubGetApi(githubServer.URL + "/200")

				require.NoError(t, err)
				assert.Equal(t, "response from github", string(response))
			})
		})

		when("github returns a non-2XX response", func() {
			it("returns an error", func() {
				_, err := client.ExecuteGithubGetApi(githubServer.URL + "/400")

				assert.Error(t, err)
				assert.Equal(t, "error: status code: 400, response: error message from github", err.Error())
			})
		})
	})
}
