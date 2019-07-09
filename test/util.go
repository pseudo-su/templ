package test

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"testing"

	"github.com/pseudo-su/templ"
	"github.com/pseudo-su/templ/tree"
	"gotest.tools/assert"
)

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(msg, err)
	}
}

func assertContains(shouldContain string) func(*testing.T, string) {
	re := regexp.MustCompile(`(` + regexp.QuoteMeta(shouldContain) + `)`)

	return func(t *testing.T, treeString string) {
		instances := len(re.FindAllString(treeString, -1))
		testDescription := fmt.Sprintf("Evaluated template should contain `%s`", shouldContain)
		assert.Equal(t, instances > 0, true, testDescription)
	}
}

// nolint:unparam
func assertContainsN(shouldContain string, desiredCount int) func(*testing.T, string) {
	re := regexp.MustCompile(`(` + regexp.QuoteMeta(shouldContain) + `)`)

	return func(t *testing.T, treeString string) {
		instances := len(re.FindAllString(treeString, -1))
		testDescription := fmt.Sprintf("Evaluated template should contain %v instances of `%s`", desiredCount, shouldContain)
		assert.Equal(t, instances == desiredCount, true, testDescription)
	}
}

// nolint:unparam
func evalAnd(t *testing.T, templateStr string, paramsArg *map[string]string, formatArg *tree.FileType, assertFn func(*testing.T, string)) {
	format := tree.YAML
	if formatArg != nil {
		format = *formatArg
	}
	params := map[string]string{}
	if paramsArg != nil {
		params = *paramsArg
	}

	inputTree, err := tree.ReadIntoTree([]byte(templateStr), format)
	failOnError(t, err, "read into tree")
	inputDesc, err := tree.DescribeTree(inputTree)
	failOnError(t, err, "describe tree")

	t.Log("INPUT TREE:")
	t.Log(inputDesc)

	outputTree, err := templ.New().Params(params).Tree("test", inputTree).Execute()
	failOnError(t, err, "execute templ")
	outputDesc, err := tree.DescribeTree(outputTree)
	failOnError(t, err, "describe outputTree")

	t.Log("TREE DESC:")
	t.Log(outputDesc)

	assertFn(t, outputDesc)
}
