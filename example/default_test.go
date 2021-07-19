package example_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	testspace "gitlab.alibaba-inc.com/agit/gotestspace"
)

// The sample test for running "echo hello"
func TestSampleShell(t *testing.T) {
	assert := assert.New(t)
	workspace, err := testspace.Create(testspace.WithShellOption("echo hello"))
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%v", err)
	}
	defer workspace.Cleanup()

	assert.Equal("hello", strings.TrimSpace(workspace.GetOutputStr()))
}

// Add environment example
func TestSampleShellWithEnvironments(t *testing.T) {
	assert := assert.New(t)
	workspace, err := testspace.Create(
		// Add two environments Testing1 and Testing2
		testspace.WithEnvironmentsOption("Testing1=aa", "Testing2=bb"),
		testspace.WithShellOption("echo $Testing1, $Testing2"),
	)
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%v", err)
	}
	defer workspace.Cleanup()

	assert.Equal("aa, bb", strings.TrimSpace(workspace.GetOutputStr()))
}

// Add template example
func TestAddTemplateAndCall(t *testing.T) {
	assert := assert.New(t)
	workspace, err := testspace.Create(
		testspace.WithTemplateOption(`
test(){
	echo "this is a test from test method"
}
`),
		testspace.WithShellOption("test"))
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%v", err)
	}
	defer workspace.Cleanup()

	assert.Equal("this is a test from test method", strings.TrimSpace(workspace.GetOutputStr()))
}

// Add custom path example
func TestSetCustomPathForTesting(t *testing.T) {
	assert := assert.New(t)
	currentPath, _ := os.Getwd()
	testFolderName := "testing_folder"
	testPath := path.Join(currentPath, testFolderName)
	workspace, err := testspace.Create(
		testspace.WithPathOption(testPath),
	)
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%v", err)
	}
	defer workspace.Cleanup()

	workPath := workspace.GetPath(testFolderName)

	_, err = os.Stat(workPath)
	if err != nil {
		assert.Error(err, "the custom path not exist")
	}
}

// Create a bare repository example
func TestCreateBareRepository(t *testing.T) {
	assert := assert.New(t)

	// The "test_tick" is the default method in template
	workspace, err := testspace.Create(
		testspace.WithShellOption(`
git init --bare test.git &&
git clone test.git test && 
(
	cd test && 
	echo "this is a test">init.js &&
	git add init.js &&
	test_tick &&
	git commit -m "this is the first commit" &&
	git push 
) &&
rm -rf test
`))
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%s", workspace.GetOutputStr())
	}
	defer workspace.Cleanup()

	// Let's add the second commit, running custom shell again
	_, _, err = workspace.Execute(`
git clone test.git test && 
(
	cd test && 
	echo "add a new file"> main.go && 
	git add main.go && 
	test_tick &&
	git commit -m "this is the second commit" && 
	git push 
)&&
rm -rf test
`)
	if !assert.NoError(err, "create testspace got error") {
		assert.FailNowf("create testspace got error", "%s", workspace.GetOutputStr())
	}

	// Now, let's check the bare repository
	_, _, err = workspace.Execute(`cd test.git && git log --oneline`)
	if !assert.NoError(err) {
		assert.FailNowf("create testspace got error", "%v", err)
	}

	assert.Equal("95dbed8 this is the second commit\n5a1f64b this is the first commit",
		strings.TrimSpace(workspace.GetOutputStr()))
}
