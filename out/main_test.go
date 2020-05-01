package main_test

import (
	"fmt"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var outPath string

func TestOut(t *testing.T) {
	buildOut(t)
	defer os.Remove(outPath)

	spec.Run(t, "Out", testOut, spec.Report(report.Terminal{}))
}

func testOut(t *testing.T, when spec.G, it spec.S) {
}

func buildOut(t *testing.T) {
	var err error
	if _, err = os.Stat("/opt/resource/out"); err == nil {
		outPath = "/opt/resource/out"
	} else {
		outFile, err := ioutil.TempFile("", "create-pull-request-resource")
		require.NoError(t, err)
		outPath = outFile.Name()
		err = outFile.Close()
		require.NoError(t, err)

		buildOutput, err := exec.Command("go", "build", "-mod", "vendor", "-o", outPath, "github.com//out").CombinedOutput()
		require.NoError(t, err, string(buildOutput))
	}
}

func runOut(testDir, s3Prefix, dependencyFile, metadataFile string) (string, error) {
	command := exec.Command(outPath, testDir)
	stdin := fmt.Sprintf(`
{
  "source": {
    "github_token": "%s",
    "repository": "%s",
    "base": "%s",
    "description":"%s",
    "branch_prefix": "%s",
    "location": "%s",
    "auto_merge": "%s",
  },
  "params": {}

}`,
	"abcv",
	"test/test",
	"master",
	"this is s PR",
	"tmp/abc",
	true,
	)
	command.Stdin = strings.NewReader(stdin)
	output, err := command.CombinedOutput()
	return string(output), err
}