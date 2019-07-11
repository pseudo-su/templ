package tree

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

type ContainerNode interface {
	forEach(func(childNode NodeRef, idx int, desc ChildKey) error) error
	getChild(ChildKey) (*NodeRef, error)
	count() int
}

type ChildKey string

// nolint:unparam
func (k *ChildKey) asMapKey() (string, error) {
	str := *k
	// lastIdx := len(str) - 1
	// if string(str[0]) == "[" && string(str[lastIdx]) == "]" {
	// 	return string(str[1:lastIdx]), nil
	// }
	return string(str), nil
}

func (k *ChildKey) asArrayIdx() (int, error) {
	str := *k
	lastIdx := len(str) - 1
	if string(str[0]) == "[" && string(str[lastIdx]) == "]" && lastIdx > 1 {
		intStr := string(str[1:lastIdx])
		result, err := strconv.Atoi(intStr)
		if err != nil {
			return -1, err
		}
		return result, nil
	}
	return -1, errors.New("invalid key")
}

type ObjectNode struct {
	raw        map[string]NodeRef
	sortedKeys []string
}

func (n *ObjectNode) nodeType() NodeType {
	return TObjectNode
}

func (n *ObjectNode) count() int {
	return len(n.sortedKeys)
}

func (n *ObjectNode) getChild(childKey ChildKey) (*NodeRef, error) {
	key, err := childKey.asMapKey()
	if err != nil {
		return nil, err
	}
	result := n.raw[key]
	return &result, nil
}

func (n *ObjectNode) isContainer() bool {
	return true
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

func (n *ObjectNode) forEach(fn func(node NodeRef, idx int, childKey ChildKey) error) error {
	for idx, key := range n.sortedKeys {
		childNode := n.raw[key]
		childKey := ChildKey(key)
		err := fn(childNode, idx, childKey)
		if err != nil {
			return err
		}
	}
	return nil
}

type ArrayNode struct {
	raw []NodeRef
}

func (n *ArrayNode) nodeType() NodeType {
	return TArrayNode
}

func (n *ArrayNode) count() int {
	return len(n.raw)
}

func (n *ArrayNode) isContainer() bool {
	return true
}

func (n *ArrayNode) getChild(childKey ChildKey) (*NodeRef, error) {
	key, err := childKey.asArrayIdx()
	if err != nil {
		return nil, err
	}
	result := n.raw[key]
	return &result, nil
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

func (n *ArrayNode) forEach(fn func(node NodeRef, idx int, childKey ChildKey) error) error {
	for idx := range n.raw {
		childNode := n.raw[idx]
		childKey := ChildKey(fmt.Sprintf("%v", idx))
		err := fn(childNode, idx, childKey)
		if err != nil {
			return err
		}
	}
	return nil
}
