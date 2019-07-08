package tree

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type ResolveContext struct {
	prevPaths [][]ChildKey
	path      []ChildKey
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
		case TStringNode, TNumberNode, TBoolNode, TNullNode:
			*resolvedStringableNodes = append(*resolvedStringableNodes, nodeRef)
			node := nodeRef.node().(StringableNode)
			return node.valStr(), nil
		case TObjectNode, TArrayNode:
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
	// var node Node = &NumberNode{raw: float64(999)}
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

	// TODO: check if there is a single resolveStringableNodes
	// if it exactly matches the
	if len(resolvedStringableNodes) == 1 {
		node := resolvedStringableNodes[0].node().(StringableNode)
		valStr := node.valStr()
		if strings.TrimSpace(result) == valStr {
			rctx.curr.swapNode(resolvedStringableNodes[0].node())
			return nil
		}
	}

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
				prevPaths: [][]ChildKey{},
				path:      ctx.path,
				curr:      ctx.curr,
				root:      ctx.root,
				parent:    ctx.parent,
			}
			return ResolveStringNode(rootTemplate, rctx)
		case TObjectNode, TArrayNode, TBoolNode, TNullNode, TNumberNode:
			return nil
		default:
			return errors.New("unknown node type")
		}
	})

	if err != nil {
		return nil, err
	}

	return rootNode, nil
}
