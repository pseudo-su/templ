package tree

import (
	"errors"
)

type WalkingContext struct {
	path   []ChildKey
	curr   NodeRef
	root   NodeRef
	parent *NodeRef
}

type NodeWalkFn func(WalkingContext) error

func WalkTree(node NodeRef, preFn NodeWalkFn) error {
	ctx := WalkingContext{
		curr:   node,
		root:   node,
		parent: nil,
	}
	return WalkTreeCtx(ctx, preFn)
}

func WalkTreeCtx(ctx WalkingContext, preFn NodeWalkFn) error {
	// Invoke action on current node
	err := preFn(ctx)
	if err != nil {
		return err
	}

	switch ctx.curr.nodeType() {
	case TObjectNode, TArrayNode:
		containerNode := ctx.curr.node().(ContainerNode)
		return containerNode.forEach(func(childNode NodeRef, idx int, childKey ChildKey) error {
			childPath := []ChildKey{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, childKey)
			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: &ctx.curr,
			}
			return WalkTreeCtx(childContext, preFn)
		})
	case TStringNode, TNumberNode, TBoolNode, TNullNode:
		return nil
	default:
		return errors.New("unknown node type found while walking")
	}
}
