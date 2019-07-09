package tree

import (
	"errors"
	"reflect"
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

type Node interface {
	nodeType() NodeType
	isContainer() bool
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
