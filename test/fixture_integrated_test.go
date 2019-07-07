package test

import (
	"os"
	"testing"

	cupaloy "github.com/bradleyjkemp/cupaloy"
	"github.com/pseudo-su/templ"
	"gotest.tools/assert"
)

func TestDifferentFiletypeLoad(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	jsonTree, err := templ.New().Params(map[string]string{}).File("fixtures/integrated/index.json").Execute()
	failOnError(t, err)
	yamlTree, err := templ.New().Params(map[string]string{}).File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err)
	tomlTree, err := templ.New().Params(map[string]string{}).File("fixtures/integrated/index.toml").Execute()
	failOnError(t, err)

	one, err := templ.DescribeTree(jsonTree, nil)
	failOnError(t, err)
	two, err := templ.DescribeTree(yamlTree, nil)
	failOnError(t, err)
	three, err := templ.DescribeTree(tomlTree, nil)
	failOnError(t, err)

	t.Log("JSON")
	t.Log(one)
	t.Log("YAML")
	t.Log(two)
	t.Log("TOML")
	t.Log(three)

	// Assert all trees are the same
	assert.DeepEqual(t, one, two)
	assert.DeepEqual(t, two, three)
}

func TestSnapshotTree(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")

	tree, err := templ.New().Params(map[string]string{}).File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err)

	cupaloy.SnapshotT(t, tree)
}

func TestSnapshotDescription(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")

	tree, err := templ.New().Params(map[string]string{}).File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err)

	treeDescription, err := templ.DescribeTree(tree, nil)
	failOnError(t, err)

	cupaloy.SnapshotT(t, treeDescription)
}
