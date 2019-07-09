package tree

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func Test_ValidNewNodePaths(t *testing.T) {
	validPathStrs := []string{
		"root",
		"[0]",
		"root.child",
		"root[child]",
		"root[0].child",
		"root[0][child]",
		"root[0].child[0].hello.world",
		// TODO: allow quotes in array notation selector
		// "root[\"hello\"]",
		// "root[\"hello]\"]",
		// "root[\"[hello\"]",
	}

	paths := []NodePath{}
	for _, str := range validPathStrs {
		path, err := NewNodePath(str)
		failOnError(t, err, "")
		paths = append(paths, *path)
	}

	cupaloy.SnapshotT(t, paths)
}
func Test_InalidNewNodePaths(t *testing.T) {
	invalidPathStrs := []string{
		"",
		// TODO
		// "root[[]",
		"root[]",
		"root.[0]",
		"root.",
	}
	errs := []error{}
	for _, str := range invalidPathStrs {
		_, err := NewNodePath(str)
		errs = append(errs, err)
	}

	cupaloy.SnapshotT(t, errs)
}
