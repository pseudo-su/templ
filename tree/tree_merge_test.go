package tree

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func Test_MergeObject(t *testing.T) {
	baseTree, err := ReadIntoTree([]byte(`
bool: true
bool_override: true
string: string to keep
string_override: string to be overwritten
`), YAML)
	failOnError(t, err, "load base tree")

	overlayTree, err := ReadIntoTree([]byte(`
bool_override: false
string_override: overwritten string
new_string: new string
`), YAML)
	failOnError(t, err, "load overlay tree")

	resultTree, err := mergeNodes(baseTree, overlayTree)

	failOnError(t, err, "merge trees")

	resultDesc, err := DescribeTree(resultTree)
	failOnError(t, err, "describe tree")

	err = cupaloy.SnapshotMulti("result tree", resultTree)
	failOnError(t, err, "snapshot")
	err = cupaloy.SnapshotMulti("result desc", resultDesc)
	failOnError(t, err, "snapshot")
}

func Test_MergeArray(t *testing.T) {
	baseTree, err := ReadIntoTree([]byte(`[ARRAY, OF, STRINGS]`), YAML)
	failOnError(t, err, "load base tree")

	overlayTree, err := ReadIntoTree([]byte(`[SOME, MORE, STRINGS]`), YAML)
	failOnError(t, err, "load overlay tree")

	resultTree, err := mergeNodes(baseTree, overlayTree)
	failOnError(t, err, "merge trees")

	resultDesc, err := DescribeTree(resultTree)
	failOnError(t, err, "describe tree")

	cupaloy.SnapshotT(t, resultDesc)
}
