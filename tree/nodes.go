package tree

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type NodeType int

const (
	TObjectNode NodeType = iota
	TArrayNode
	TNumberNode
	TStringNode
	TBoolNode
	TNullNode
)

type NodeRef interface {
	swapNode(n Node)
	nodeType() NodeType
	// raw() interface{}
	node() Node
	forEachChild(func(childNode NodeRef, idx int, desc string) error) error
}

type NodeReference struct {
	n Node
}

func (ref *NodeReference) node() Node {
	return ref.n
}

func (ref *NodeReference) swapNode(n Node) {
	ref.n = n
}

func (ref *NodeReference) nodeType() NodeType {
	return ref.n.nodeType()
}

func (ref *NodeReference) forEachChild(executeFn func(childNode NodeRef, idx int, desc string) error) error {
	switch node := ref.n.(type) {
	case *ObjectNode:
		return node.forEach(executeFn)
	case *ArrayNode:
		return node.forEach(executeFn)
	}
	return nil
}

type Node interface {
	nodeType() NodeType
}

func NewNode(v reflect.Value) (NodeRef, error) {
	switch v.Kind() {
	case reflect.Map:
		return NewObjectNodeRef(v)
	case reflect.Slice:
		return NewArrayNodeRef(v)
	case reflect.Int64:
		var val = v.Interface().(int64)
		floatVal := float64(val)
		var node Node = &NumberNode{raw: floatVal}
		ref := &NodeReference{n: node}
		return ref, nil
	case reflect.Int:
		var val = v.Interface().(int)
		floatVal := float64(val)
		node := &NumberNode{raw: floatVal}
		ref := &NodeReference{n: node}
		return ref, nil
	case reflect.Float64:
		var val = v.Interface().(float64)
		node := &NumberNode{raw: val}
		ref := &NodeReference{n: node}
		return ref, nil
	case reflect.Bool:
		var val = v.Interface().(bool)
		node := &BoolNode{raw: val}
		ref := &NodeReference{n: node}
		return ref, nil
	case reflect.String:
		var val = v.Interface().(string)
		node := &StringNode{raw: val, templateResolved: false}
		ref := &NodeReference{n: node}
		return ref, nil
	default:
		return nil, errors.New("unable to build tree")
	}
}

type ObjectNode struct {
	raw        map[string]NodeRef
	sortedKeys []string
}

func (n *ObjectNode) nodeType() NodeType {
	return TObjectNode
}

type Entry struct {
	k string
	v reflect.Value
}

func NewObjectNodeRef(v reflect.Value) (NodeRef, error) {
	node, err := NewObjectNode(v)
	if err != nil {
		return nil, err
	}
	ref := &NodeReference{n: node}
	return ref, nil
}

func NewObjectNode(v reflect.Value) (*ObjectNode, error) {
	iter := v.MapRange()
	var entries = []Entry{}
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		entries = append(entries, Entry{
			k: k.Interface().(string),
			v: v,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].k < entries[j].k
	})
	var raw = make(map[string]NodeRef)
	sortedKeys := []string{}
	for i := 0; i < len(entries); i++ {
		e := entries[i]
		nodeRef, err := NewNode(e.v.Elem())
		if err != nil {
			return nil, err
		}
		raw[e.k] = nodeRef
		sortedKeys = append(sortedKeys, e.k)
	}
	node := ObjectNode{raw: raw, sortedKeys: sortedKeys}
	return &node, nil
}

func (n *ObjectNode) forEach(fn func(node NodeRef, idx int, desc string) error) error {
	for idx, key := range n.sortedKeys {
		childNode := n.raw[key]
		desc := key
		err := fn(childNode, idx, desc)
		if err != nil {
			return err
		}
	}
	return nil
}

type StringNode struct {
	raw              string
	templateResolved bool
}

func (n *StringNode) nodeType() NodeType {
	return TStringNode
}

func (n *StringNode) updateNode(newVal interface{}) error {
	switch val := newVal.(type) {
	case string:
		n.raw = val
		return nil
	case *StringNode:
		n.raw = val.raw
		return nil
	case StringNode:
		n.raw = val.raw
		return nil
	default:
		return errors.New("unable to set new string value")
	}
}

type NumberNode struct {
	raw float64
}

func (n *NumberNode) nodeType() NodeType {
	return TNumberNode
}

type ArrayNode struct {
	raw []NodeRef
}

func (n *ArrayNode) nodeType() NodeType {
	return TArrayNode
}

func NewArrayNodeRef(v reflect.Value) (NodeRef, error) {
	node, err := NewArrayNode(v)
	if err != nil {
		return nil, err
	}
	ref := &NodeReference{n: node}
	return ref, nil
}

func NewArrayNode(v reflect.Value) (*ArrayNode, error) {
	var raw = []NodeRef{}

	for i := 0; i < v.Len(); i++ {
		n, err := NewNode(v.Index(i).Elem())
		if err != nil {
			return &ArrayNode{}, err
		}
		raw = append(raw, n)
	}

	node := ArrayNode{raw: raw}

	return &node, nil
}

func (n *ArrayNode) forEach(fn func(node NodeRef, idx int, desc string) error) error {
	for idx, _ := range n.raw {
		childNode := n.raw[idx]
		desc := fmt.Sprintf("[%v]", idx)
		err := fn(childNode, idx, desc)
		if err != nil {
			return err
		}
	}
	return nil
}

type BoolNode struct {
	raw bool
}

func (n *BoolNode) nodeType() NodeType {
	return TBoolNode
}

type NullNode struct{}

func (n *NullNode) nodeType() NodeType {
	return TNullNode
}
