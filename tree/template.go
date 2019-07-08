package tree

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type ResolveContext struct {
	prevPaths [][]string
	path      []string
	curr      *Node
	root      *Node
	parent    *Node
}

// func ResolveInTemplate

func CreateSelfFn(rctx ResolveContext, resolvedNodes *[]Node) func(...string) (string, error) {
	return func(args ...string) (string, error) {
		// when resolving an object replace
		node := *ResolvePath(JoinStr(args), rctx)
		switch node := node.(type) {
		case StringNode:
			return node.raw, nil
		case NumberNode:
			return fmt.Sprintf("%v", node.raw), nil
		case BoolNode:
			return fmt.Sprintf("%v", node.raw), nil
		case ObjectNode:
			*resolvedNodes = append(*resolvedNodes, node)
			return "", nil
		case ArrayNode:
			*resolvedNodes = append(*resolvedNodes, node)
			return "", nil
		default:
			return "", errors.New("unhandled Node type")
		}
	}
}

func CreateFileLoadFn(rctx ResolveContext, resolvedNodes *[]Node) func(...string) (string, error) {
	return func(args ...string) (string, error) {
		// when resolving an object replace
		return "<TODO: load>", nil
	}
}

func ResolvePath(path string, rctx ResolveContext) *Node {
	// TODO: resolve to the node at a given path
	var node Node = StringNode{raw: "test_value_todo", templateResolved: true}
	return &node
}

func ResolveStringNode(rootTemplate template.Template, rctx ResolveContext) error {
	val := (*rctx.curr).(StringNode).raw
	// resolvedNodes := map[string]interface{}{}
	resolvedNodes := []Node{}
	tmpl, err := rootTemplate.Funcs(template.FuncMap{
		"self": CreateSelfFn(rctx, &resolvedNodes),
		"load": CreateFileLoadFn(rctx, &resolvedNodes),
	}).Parse(val)
	if err != nil {
		return err
	}
	// TODO pass params to template
	data := struct {
		stage string
	}{
		stage: "Some string",
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	result := tpl.String()

	if len(resolvedNodes) > 1 {
		return errors.New("a single field can not resolve to multiple complex values")
	}

	if len(resolvedNodes) == 1 && strings.TrimSpace(result) != "" {
		return errors.New("unable to join strings with complex values")
	}

	if len(resolvedNodes) == 1 {
		fmt.Println("replace with node")
		*rctx.curr = resolvedNodes[0]
	} else {
		fmt.Println("update string node:", result)
		*rctx.curr = StringNode{raw: result, templateResolved: true}

		// TODO: this
		fmt.Println("from parent:", (*rctx.parent).(ObjectNode).raw[rctx.path[len(rctx.path)-1]].(StringNode).raw)
		fmt.Println("from current:", (*rctx.curr).(StringNode).raw)
	}

	return nil
}

func ExecuteTreeTemplate(rootNode *Node, rootTemplate template.Template) (*Node, error) {
	err := WalkTree(rootNode, func(ctx WalkingContext) error {
		node := *ctx.curr
		switch node.(type) {
		case StringNode:
			rctx := ResolveContext{
				prevPaths: [][]string{},
				path:      ctx.path,
				curr:      ctx.curr,
				root:      ctx.root,
				parent:    ctx.parent,
			}
			return ResolveStringNode(rootTemplate, rctx)
		case ObjectNode, ArrayNode:
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return rootNode, nil
}
