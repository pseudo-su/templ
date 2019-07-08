package tree

import (
	"errors"
)

type WalkingContext struct {
	path   []string
	curr   NodeRef
	root   NodeRef
	parent *NodeRef
}

type NodeWalkFn func(WalkingContext) error

func WalkTree(node NodeRef, actionFn NodeWalkFn) error {
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

	switch ctx.curr.nodeType() {
	case TObjectNode:
		return ctx.curr.forEachChild(func(childNode NodeRef, idx int, desc string) error {
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, desc)
			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: &ctx.curr,
			}
			return WalkTreeCtx(childContext, actionFn)
		})
	case TArrayNode:
		return ctx.curr.forEachChild(func(childNode NodeRef, idx int, desc string) error {
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, desc)

			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: &ctx.curr,
			}
			return WalkTreeCtx(childContext, actionFn)
		})
	case TStringNode, TNumberNode, TBoolNode, TNullNode:
		return nil
	default:
		return errors.New("unknown node type found while walking")
	}
}
