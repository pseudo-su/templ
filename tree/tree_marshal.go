package tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// TODO: implement unmarshaling into golang interface/types
// func UnmarshalTree(tree NodeRef, v interface{}) error {
// 	return nil
// }

// TODO: get rid of HTML encoding
func jsonEscape(i string) (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(b[1 : len(b)-1]), nil
}

// TODO: implement JSON marshaling
func (n *NodeReference) MarshalJSON() ([]byte, error) {
	ctx := WalkingContext{
		curr:   n,
		root:   n,
		parent: nil,
	}
	return marshalJSON(ctx)
}

// nolint:gocyclo
func marshalJSON(ctx WalkingContext) ([]byte, error) {
	result := ""
	switch ctx.curr.nodeType() {
	case TObjectNode:
		containerNode := ctx.curr.node().(*ObjectNode)
		result += "{"
		err := containerNode.forEach(func(childNode NodeRef, idx int, childKey ChildKey) error {
			childPath := []ChildKey{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, childKey)
			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: &ctx.curr,
			}
			str, err := marshalJSON(childContext)
			if err != nil {
				return err
			}
			key, err := childKey.asMapKey()
			if err != nil {
				return err
			}
			result += fmt.Sprintf("\"%v\":", key)
			result += string(str)
			if idx != containerNode.count()-1 {
				result += ","
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		result += "}"
	case TArrayNode:
		containerNode := ctx.curr.node().(*ArrayNode)
		result += "["
		err := containerNode.forEach(func(childNode NodeRef, idx int, childKey ChildKey) error {
			childPath := []ChildKey{}
			childPath = append(childPath, ctx.path...)
			childPath = append(childPath, childKey)
			childContext := WalkingContext{
				path:   childPath,
				root:   ctx.root,
				curr:   childNode,
				parent: &ctx.curr,
			}
			str, err := marshalJSON(childContext)
			if err != nil {
				return err
			}
			result += string(str)
			if idx != containerNode.count()-1 {
				result += ","
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		result += "]"
	case TStringNode:
		strNode := ctx.curr.node().(StringableNode)
		val, err := jsonEscape(strNode.valStr())
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf("\"%v\"", val)
	case TNumberNode, TBoolNode, TNullNode:
		strNode := ctx.curr.node().(StringableNode)
		result += strNode.valStr()
	default:
		return nil, errors.New("unknown node type found while walking")
	}
	return []byte(result), nil
}

func (n *NodeReference) MarshalJSONIndent(prefix, indent string) ([]byte, error) {
	b, err := n.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// TODO: implement YAML marshaling
// func (n *NodeReference) MarshalYAML() {

// }
