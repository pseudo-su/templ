package tree

import (
	// "errors"

	"fmt"
	"strings"
)

func descPrefix(ctx WalkingContext) (string, string) {
	depth := len(ctx.path)
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "  "
	}
	if prefix == "" {
		prefix = "root"
	}
	label := ""
	if depth > 0 {
		label = ctx.path[depth-1]
	}
	return prefix, label
}

func DescribeTree(rootNode Node) (string, error) {
	lines := []string{}
	err := WalkTree(rootNode, func(ctx WalkingContext) error {
		indent, label := descPrefix(ctx)
		switch n := ctx.curr.(type) {
		case *ObjectNode:
			lines = append(lines, Indent(fmt.Sprintf("%v(object):", label), indent))
		case *ArrayNode:
			lines = append(lines, Indent(fmt.Sprintf("%v(array):", label), indent))
		case *StringNode:
			val := *n.raw
			if strings.Contains(val, "\n") {
				lines = append(lines, Indent(fmt.Sprintf("%v(mstr):", label), indent))
				lines = append(lines, Indent(val, indent+" > "))
			} else {
				lines = append(lines, Indent(fmt.Sprintf("%v(str): %v", label, val), indent))
			}
		case *NumberNode:
			val := *n.raw
			lines = append(lines, fmt.Sprintf("%v(num): %v", label, val))
		case *BoolNode:
			val := *n.raw
			lines = append(lines, fmt.Sprintf("%v(bool): %v", label, val))
		case *NullNode:
			lines = append(lines, fmt.Sprintf("%v: null", label))
		default:
			return fmt.Errorf("unable to describe node")
		}
		return nil
	})
	return strings.Join(lines, "\n"), err
}

func Indent(s, prefix string) string {
	return string(IndentBytes([]byte(s), []byte(prefix)))
}

func IndentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}
	// Also indent after trailing newline
	if bol {
		res = append(res, prefix...)
	}
	return res
}
