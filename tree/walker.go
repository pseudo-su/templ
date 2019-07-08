package tree

import (
	"errors"
)

// import "errors"

type WalkingContext struct {
	path   []string
	curr   *Node
	root   *Node
	parent *Node
}

type NodeWalkFn func(WalkingContext) error

func WalkTree(node *Node, actionFn NodeWalkFn) error {
	ctx := WalkingContext{
		curr:   node,
		root:   node,
		parent: nil,
	}
	return WalkTreeCtx(ctx, actionFn)
}

func WalkTreeCtx(ctx WalkingContext, actionFn NodeWalkFn) error {
	// Invoke action on current node
	err := actionFn(ctx)
	if err != nil {
		return err
	}

	// Continue walking if applicable
	node := *ctx.curr

	switch node := node.(type) {
	case ObjectNode:
		return node.forEach(func(childNode *Node, idx int, desc string) error {
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, desc)
			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: ctx.curr,
			}
			return WalkTreeCtx(childContext, actionFn)
		})
	case ArrayNode:
		return node.forEach(func(childNode *Node, idx int, desc string) error {
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, desc)

			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: ctx.curr,
			}
			return WalkTreeCtx(childContext, actionFn)
		})
	case StringNode, NumberNode, BoolNode, NullNode:
		return nil
	default:
		return errors.New("unknown node type found while walking")
	}
}
