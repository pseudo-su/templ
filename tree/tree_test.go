package tree

import (
	"testing"

	"gotest.tools/assert"
)

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatal(msg, err)
	}
}

// func TestDescribeTree(t *testing.T) {
// 	tree, err := ReadIntoTree([]byte(`
// ["array", "of", "strings"]
// `), JSON)
// 	failOnError(t, err, "")
// 	treeDesc, err := DescribeTree(tree)
// 	failOnError(t, err, "")
// }

func TestTreesMatch(t *testing.T) {
	jsonTree, err := ReadIntoTree([]byte(`
{
	"array": ["array", "of", "strings"],
	"bool": true,
	"float": 1.1,
	"int": 1,
	"object": { "a": "a" },
	"string": "string"
}
`), JSON)
	failOnError(t, err, "")
	yamlTree, err := ReadIntoTree([]byte(`
int: 1
float: 1.1
string: string
bool: true
array:
  - array
  - of
  - strings
object:
  a: a
`), YAML)
	failOnError(t, err, "")
	tomlTree, err := ReadIntoTree([]byte(`
array = ["array", "of", "strings"]
bool = true
float = 1.1
int = 1
string = "string"

[object]
a = "a"
`), TOML)
	failOnError(t, err, "")

	one, err := DescribeTree(jsonTree)
	failOnError(t, err, "DescribeTree (json)")
	two, err := DescribeTree(yamlTree)
	failOnError(t, err, "DescribeTree (yaml)")
	three, err := DescribeTree(tomlTree)
	failOnError(t, err, "DescribeTree (toml)")

	t.Log("JSON")
	t.Log(one)
	t.Log("YAML")
	t.Log(two)
	t.Log("TOML")
	t.Log(three)

	assert.DeepEqual(t, one, two)
	assert.DeepEqual(t, two, three)
}
