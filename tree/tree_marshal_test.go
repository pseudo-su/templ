package tree

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func Test_Marshal(t *testing.T) {
	testTree, err := ReadIntoTree([]byte(`
bool: true
bool_override: true
string: string to keep
string_override: string to be overwritten
multi_string: |
  hello world
  this is me
object:
  some: val
  nested:
    val: another
array:
  - one
  - two
  - { a: [a, b, c, null, 1, 2.1] }
`), YAML)
	failOnError(t, err, "load base tree")

	// Normal JSON stringify
	resultJSON, err := testTree.MarshalJSON()
	failOnError(t, err, "marshal")
	err = cupaloy.SnapshotMulti("simple JSON", resultJSON)
	failOnError(t, err, "snapshot simple")

	// JSON stringify pretty
	prettyJSON, err := testTree.MarshalJSONIndent("", "  ")
	failOnError(t, err, "marshal")
	err = cupaloy.SnapshotMulti("pretty JSON", prettyJSON)
	failOnError(t, err, "snapshot pretty")
}
