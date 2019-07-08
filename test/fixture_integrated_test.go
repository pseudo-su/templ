package test

import (
	"os"
	"testing"

	cupaloy "github.com/bradleyjkemp/cupaloy"
	"github.com/pseudo-su/templ"
	templTree "github.com/pseudo-su/templ/tree"
	"gotest.tools/assert"
)

func TestDifferentFiletypeLoad(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")
	jsonTree, err := templ.New().File("fixtures/integrated/index.json").Execute()
	failOnError(t, err, "load index.json")
	yamlTree, err := templ.New().File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err, "load index.yaml")
	tomlTree, err := templ.New().File("fixtures/integrated/index.toml").Execute()
	failOnError(t, err, "load index.toml")

	one, err := templTree.DescribeTree(jsonTree)
	failOnError(t, err, "describe json")
	two, err := templTree.DescribeTree(yamlTree)
	failOnError(t, err, "describe yaml")
	three, err := templTree.DescribeTree(tomlTree)
	failOnError(t, err, "describe toml")

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

	tree, err := templ.New().File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err, "execute templ")

	cupaloy.SnapshotT(t, tree)
}

func TestSnapshotDescription(t *testing.T) {
	os.Setenv("ENV", "envname")
	os.Setenv("STAGE", "stagename")

	tree, err := templ.New().File("fixtures/integrated/index.yaml").Execute()
	failOnError(t, err, "execute templ")

	treeDescription, err := templTree.DescribeTree(tree)
	failOnError(t, err, "describe tree")

	cupaloy.SnapshotT(t, treeDescription)
}
