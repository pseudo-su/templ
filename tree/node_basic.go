package tree

import (
	"fmt"
)

type StringableNode interface {
	valStr() string
}

type StringNode struct {
	raw              string
	templateResolved bool
}

func (n *StringNode) nodeType() NodeType {
	return TStringNode
}

func (n *StringNode) isContainer() bool {
	return false
}

func (n *StringNode) valStr() string {
	return n.raw
}

// func (n *StringNode) updateNode(newVal interface{}) error {
// 	switch val := newVal.(type) {
// 	case string:
// 		n.raw = val
// 		return nil
// 	case *StringNode:
// 		n.raw = val.raw
// 		return nil
// 	case StringNode:
// 		n.raw = val.raw
// 		return nil
// 	default:
// 		return errors.New("unable to set new string value")
// 	}
// }

type NumberNode struct {
	raw float64
}

func (n *NumberNode) nodeType() NodeType {
	return TNumberNode
}

func (n *NumberNode) isContainer() bool {
	return false
}

func (n *NumberNode) valStr() string {
	return fmt.Sprintf("%v", n.raw)
}

type BoolNode struct {
	raw bool
}

func (n *BoolNode) nodeType() NodeType {
	return TBoolNode
}

func (n *BoolNode) isContainer() bool {
	return false
}

func (n *BoolNode) valStr() string {
	return fmt.Sprintf("%v", n.raw)
}

type NullNode struct{}

func (n *NullNode) nodeType() NodeType {
	return TNullNode
}

func (n *NullNode) isContainer() bool {
	return false
}

func (n *NullNode) valStr() string {
	return "null"
}
