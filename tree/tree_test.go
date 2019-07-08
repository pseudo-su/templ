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

func TestFloatRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`2.11`), YAML)
	failOnError(t, err, "parse float")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc float")
	t.Log(desc)
	assert.DeepEqual(t, desc, "(num): 2.11")
}

func TestIntegerRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`2`), YAML)
	failOnError(t, err, "parse int")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc int")
	t.Log(desc)
	assert.DeepEqual(t, desc, "(num): 2")
}

func TestStringRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`"hello world"`), YAML)
	failOnError(t, err, "parse str")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc str")
	t.Log(desc)
	assert.DeepEqual(t, desc, "root(str): hello world")
}

func TestBoolRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`false`), YAML)
	failOnError(t, err, "parse bool")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc bool")
	t.Log(desc)
	assert.DeepEqual(t, desc, "(bool): false")
}

func TestArrayRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`[MY, WORD]`), YAML)
	failOnError(t, err, "parse array")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc array")
	t.Log(desc)
	assert.DeepEqual(t, desc, "root(array):\n  [0](str): MY\n  [1](str): WORD")
}

func TestObjectRoot(t *testing.T) {
	tree, err := ReadIntoTree([]byte(`{"key":"val"}`), JSON)
	failOnError(t, err, "parse obj")
	desc, err := DescribeTree(tree)
	failOnError(t, err, "desc obj")
	t.Log(desc)
	assert.DeepEqual(t, desc, "root(object):\n  key(str): val")
}

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
	failOnError(t, err, "load json")
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
	failOnError(t, err, "load yaml")
	tomlTree, err := ReadIntoTree([]byte(`
array = ["array", "of", "strings"]
bool = true
float = 1.1
int = 1
string = "string"

[object]
a = "a"
`), TOML)
	failOnError(t, err, "load toml")

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
