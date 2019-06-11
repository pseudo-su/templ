package tree

import (
	"fmt"
)

type WalkingContext struct {
	path   []string
	curr   Node
	root   Node
	parent Node
}

type NodeWalkFn func(WalkingContext) error

func WalkTree(node Node, actionFn NodeWalkFn, ctx *WalkingContext) error {
	if ctx == nil {
		ctx = &WalkingContext{
			curr:   node,
			root:   node,
			parent: nil,
		}
	}
	// TODO: do something with return value?
	err := actionFn(*ctx)
	if err != nil {
		return err
	}

	switch node.(type) {
	case ObjectNode:
		node := node.(ObjectNode)

		for _, key := range node.sortedKeys {
			childNode := node.raw[key]
			// TODO: childContext
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, key)
			childContext := &WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: ctx.curr,
			}
			err := WalkTree(childNode, actionFn, childContext)
			if err != nil {
				return err
			}
		}
	case ArrayNode:
		node := node.(ArrayNode)
		for i, childNode := range node.raw {
			// TODO: childContext
			childPath := []string{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, fmt.Sprintf("[%v]", i))

			childContext := &WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: ctx.curr,
			}
			err := WalkTree(childNode, actionFn, childContext)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	return nil
}
