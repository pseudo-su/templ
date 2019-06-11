package tree

import (
	"errors"
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
	var sortedKeys = make([]string, len(entries))
	for _, e := range entries {
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

// func (n *ObjectNode) MarshalJSON() ([]byte, error) {
// 	spew.Dump("MarshalJSON")
// 	buffer := bytes.NewBufferString("{")
// 	length := len(n.raw)
// 	count := 0
// 	for key, value := range n.raw {
// 		jsonValue, err := json.Marshal(value)
// 		if err != nil {
// 			return nil, err
// 		}
// 		buffer.WriteString(fmt.Sprintf("\"%s\":%s", key, string(jsonValue)))
// 		count++
// 		if count < length {
// 			buffer.WriteString(",")
// 		}
// 	}
// 	buffer.WriteString("}")
// 	spew.Dump(buffer)
// 	return buffer.Bytes(), nil
// }

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

type BoolNode struct {
	raw bool
}

type NullNode struct{}
