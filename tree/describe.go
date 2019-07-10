package tree

import (
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
		label = fmt.Sprintf("%v", ctx.path[depth-1])
	}
	return prefix, label
}

func DescribeTree(rootNode NodeRef) (string, error) {
	lines := []string{}
	err := WalkTree(rootNode, func(ctx WalkingContext) error {
		indent, label := descPrefix(ctx)
		switch ctx.curr.nodeType() {
		case TObjectNode:
			lines = append(lines, Indent(fmt.Sprintf("%v(object):", label), indent))
		case TArrayNode:
			lines = append(lines, Indent(fmt.Sprintf("%v(array):", label), indent))
		case TStringNode:
			node := ctx.curr.node()
			val := node.(StringableNode).valStr()
			if strings.Contains(val, "\n") {
				lines = append(lines, Indent(fmt.Sprintf("%v(mstr):", label), indent))
				lines = append(lines, Indent(val, indent+" > "))
			} else {
				lines = append(lines, Indent(fmt.Sprintf("%v(str): %v", label, val), indent))
			}
		case TNumberNode:
			node := ctx.curr.node()
			val := node.(StringableNode).valStr()
			lines = append(lines, Indent(fmt.Sprintf("%v(num): %v", label, val), indent))
		case TBoolNode:
			node := ctx.curr.node()
			val := node.(StringableNode).valStr()
			lines = append(lines, Indent(fmt.Sprintf("%v(bool): %v", label, val), indent))
		case TNullNode:
			lines = append(lines, Indent(fmt.Sprintf("%v: null", label), indent))
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
