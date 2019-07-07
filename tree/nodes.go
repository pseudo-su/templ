package tree

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

type Node interface{}

func NewNode(v reflect.Value) (Node, error) {
	switch v.Kind() {
	case reflect.Map:
		return NewObjectNode(v)
	case reflect.Slice:
		return NewArrayNode(v)
	case reflect.Int64:
		var raw = v.Interface().(int64)
		node := NumberNode{raw: float64(raw)}
		return node, nil
	case reflect.Int:
		var raw = v.Interface().(int)
		node := NumberNode{raw: float64(raw)}
		return node, nil
	case reflect.Float64:
		var raw = v.Interface().(float64)
		node := NumberNode{raw: raw}
		return node, nil
	case reflect.Bool:
		var raw = v.Interface().(bool)
		node := BoolNode{raw: raw}
		return node, nil
	case reflect.String:
		var raw = v.Interface().(string)
		node := StringNode{raw: raw, templateResolved: false}
		return node, nil
	default:
		return nil, errors.New("unable to build tree")
	}
}

type ObjectNode struct {
	raw        map[string]Node
	sortedKeys []string
}
type Entry struct {
	k string
	v reflect.Value
}

func NewObjectNode(v reflect.Value) (ObjectNode, error) {
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
	var raw = make(map[string]Node)
	sortedKeys := []string{}
	for i := 0; i < len(entries); i++ {
		e := entries[i]
		val, err := NewNode(e.v.Elem())
		if err != nil {
			return ObjectNode{raw: raw, sortedKeys: sortedKeys}, err
		}
		raw[e.k] = val
		sortedKeys = append(sortedKeys, e.k)
	}
	node := ObjectNode{raw: raw, sortedKeys: sortedKeys}
	return node, nil
}

func (n *ObjectNode) forEach(fn func(node *Node, idx int, desc string) error) error {
	for idx, key := range n.sortedKeys {
		childNode := n.raw[key]
		desc := key
		err := fn(&childNode, idx, desc)
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

type NumberNode struct {
	raw float64
}

type ArrayNode struct {
	raw []Node
}

func NewArrayNode(v reflect.Value) (ArrayNode, error) {
	var raw = []Node{}

	for i := 0; i < v.Len(); i++ {
		val, err := NewNode(v.Index(i).Elem())
		if err != nil {
			return ArrayNode{}, err
		}
		raw = append(raw, val)
	}

	node := ArrayNode{raw: raw}
	// spew.Dump(node)
	return node, nil
}

func (n *ArrayNode) forEach(fn func(node *Node, idx int, desc string) error) error {
	for idx, _ := range n.raw {
		childNode := n.raw[idx]
		desc := fmt.Sprintf("[%v]", idx)
		err := fn(&childNode, idx, desc)
		if err != nil {
			return err
		}
	}
	return nil
}

type BoolNode struct {
	raw bool
}

type NullNode struct{}
