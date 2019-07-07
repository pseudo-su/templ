package test

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"testing"

	"github.com/pseudo-su/templ"
	"github.com/pseudo-su/templ/tree"
	templTree "github.com/pseudo-su/templ/tree"
	"gotest.tools/assert"
)

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(msg, err)
	}
}

func assertContains(shouldContain string) func(*testing.T, string) {
	re := regexp.MustCompile(`(` + shouldContain + `)`)

	return func(t *testing.T, treeString string) {
		instances := len(re.FindAllString(treeString, -1))
		testDescription := fmt.Sprintf("Evaluated template should contain `%s`", shouldContain)
		assert.Equal(t, instances > 0, true, testDescription)
	}
}

func assertContainsN(shouldContain string, desiredCount int) func(*testing.T, string) {
	re := regexp.MustCompile(`(` + shouldContain + `)`)

	return func(t *testing.T, treeString string) {
		instances := len(re.FindAllString(treeString, -1))
		testDescription := fmt.Sprintf("Evaluated template should contain %v instances of `%s`", desiredCount, shouldContain)
		assert.Equal(t, instances == desiredCount, true, testDescription)
	}
}

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
	failOnError(t, err, "")
	outputTree, err := templ.New().Params(params).Tree("test", inputTree).Execute()
	failOnError(t, err, "")
	treeString, err := templTree.DescribeTree(outputTree)
	failOnError(t, err, "")

	t.Log("TREE:")
	t.Log(outputTree)
	t.Log("TREE DESC:")
	t.Log(treeString)
	assertFn(t, treeString)
}
