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
	curr      NodeRef
	root      NodeRef
	parent    *NodeRef
}

// func ResolveInTemplate

func CreateSelfFn(
	rctx ResolveContext,
	resolvedStringableNodes *[]NodeRef,
	resolvedComplexNodes *[]NodeRef,
) func(...string) (string, error) {
	return func(args ...string) (string, error) {
		// when resolving an object replace
		nodeRef := ResolvePath(JoinStr(args), rctx)
		switch nodeRef.nodeType() {
		case TStringNode:
			node := nodeRef.node()
			raw := node.(*StringNode).raw
			return raw, nil
		case TNumberNode:
			// TODO: fix number resolve
			node := nodeRef.node()
			raw := node.(*NumberNode).raw
			return fmt.Sprintf("%v", raw), nil
		case TBoolNode:
			// TODO: fix bool resolve
			node := nodeRef.node()
			raw := node.(*BoolNode).raw
			return fmt.Sprintf("%v", raw), nil
		case TObjectNode:
			*resolvedComplexNodes = append(*resolvedComplexNodes, nodeRef)
			return "", nil
		case TArrayNode:
			*resolvedComplexNodes = append(*resolvedComplexNodes, nodeRef)
			return "", nil
		default:
			return "", errors.New("unhandled Node type")
		}
	}
}

func CreateFileLoadFn(
	rctx ResolveContext,
	resolvedStringableNodes *[]NodeRef,
	resolvedComplexNodes *[]NodeRef,
) func(...string) (string, error) {
	return func(args ...string) (string, error) {
		// when resolving an object replace
		return "<TODO: load>", nil
	}
}

func ResolvePath(path string, rctx ResolveContext) NodeRef {
	// TODO: resolve to the node at a given path
	var node Node = &StringNode{raw: "<TODO: self (" + path + ") >", templateResolved: true}
	// var node Node = &NumberNode{raw: float64(1)}
	// var node Node = &ArrayNode{}
	ref := &NodeReference{n: node}
	return ref
}

func ResolveStringNode(rootTemplate template.Template, rctx ResolveContext) error {
	node := rctx.curr.node()
	val := node.(*StringNode).raw
	// resolvedNodes := map[string]interface{}{}
	resolvedStringableNodes := []NodeRef{}
	resolvedComplexNodes := []NodeRef{}
	tmpl, err := rootTemplate.Funcs(template.FuncMap{
		"self": CreateSelfFn(rctx, &resolvedStringableNodes, &resolvedComplexNodes),
		"load": CreateFileLoadFn(rctx, &resolvedStringableNodes, &resolvedComplexNodes),
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

	if len(resolvedComplexNodes) > 1 {
		fmt.Println(resolvedComplexNodes)
		return errors.New("a single field can not resolve to multiple complex values")
	}

	if len(resolvedComplexNodes) == 1 && strings.TrimSpace(result) != "" {
		return errors.New("unable to join strings with complex values")
	}

	if len(resolvedComplexNodes) == 1 {
		fmt.Println("replace string node")
		node := resolvedComplexNodes[0].node()
		rctx.curr.swapNode(node)
		return nil
	}
	fmt.Println("update string node:", result)
	// TODO update string node instead of replace
	newNode := &StringNode{raw: result, templateResolved: true}
	rctx.curr.swapNode(newNode)

	return nil
}

func ExecuteTreeTemplate(rootNode NodeRef, rootTemplate template.Template) (NodeRef, error) {
	err := WalkTree(rootNode, func(ctx WalkingContext) error {
		switch ctx.curr.nodeType() {
		case TStringNode:
			rctx := ResolveContext{
				prevPaths: [][]string{},
				path:      ctx.path,
				curr:      ctx.curr,
				root:      ctx.root,
				parent:    ctx.parent,
			}
			return ResolveStringNode(rootTemplate, rctx)
		case TObjectNode, TArrayNode:
			return nil
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return rootNode, nil
}
