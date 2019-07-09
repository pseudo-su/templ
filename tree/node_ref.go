package tree

import (
	"errors"
)

type NodeRef interface {
	swapNode(n Node)
	nodeType() NodeType
	node() Node
	selectNode(selectStr string) (NodeRef, error)
	selectNodeByNodePath(selectStr NodePath) (NodeRef, error)
}

type NodePath []ChildKey

// nolint:gocyclo
func NewNodePath(pathStr string) (*NodePath, error) {
	if pathStr == "" {
		return nil, errors.New("invalid path string: empty string")
	}
	paths := NodePath{}
	bytes := []byte(pathStr)
	lastEnd := 0
	closingBracesNeeded := 0
	lastCharIdx := len(bytes) - 1
	for i, char := range bytes {
		prevCharIdx := i - 1
		needCloseBrace := closingBracesNeeded > 0
		if needCloseBrace {
			if char == ']' {
				closingBracesNeeded--
			}

			if char == ']' && closingBracesNeeded == 0 {
				val := pathStr[lastEnd:i]
				key := ChildKey(val)
				if key == "" {
					return nil, errors.New("invalid path string: unexpected character \"]\"")
				}
				paths = append(paths, key)
				lastEnd = i + 1
			}
			continue
		}
		if char == '[' && lastCharIdx == i {
			return nil, errors.New("unexpected end of string")
		}
		if char == '[' && prevCharIdx == -1 {
			lastEnd = i + 1
			closingBracesNeeded++
			continue
		}
		if char == '[' && prevCharIdx >= 0 && bytes[prevCharIdx] == ']' {
			lastEnd = i + 1
			closingBracesNeeded++
			continue
		}
		if char == '[' {
			val := pathStr[lastEnd:i]
			key := ChildKey(val)
			if key == "" {
				return nil, errors.New("invalid path string: unexpected character \"[\"")
			}
			paths = append(paths, key)

			lastEnd = i + 1
			closingBracesNeeded++
			continue
		}

		if char == '.' && lastCharIdx == i {
			return nil, errors.New("invalid path string: unexpected end of string")
		}
		if char == '.' && prevCharIdx >= 0 && bytes[prevCharIdx] == ']' {
			lastEnd = i + 1
			continue
		}
		if char == '.' {
			val := pathStr[lastEnd:i]
			key := ChildKey(val)
			paths = append(paths, key)
			lastEnd = i + 1
			continue
		}
		if i == lastCharIdx {
			val := pathStr[lastEnd : i+1]
			key := ChildKey(val)
			paths = append(paths, key)
			lastEnd = i + 1
			continue
		}
	}
	if closingBracesNeeded > 0 {
		return nil, errors.New("invalid path, no closing bracket")
	}
	return &paths, nil
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

func (ref *NodeReference) selectNode(selectStr string) (NodeRef, error) {
	path, err := NewNodePath(selectStr)
	if err != nil {
		return nil, err
	}
	sRef, err := ref.selectNodeByNodePath(*path)
	if err != nil {
		return nil, err
	}
	return sRef, nil
}

func (ref *NodeReference) selectNodeByNodePath(path NodePath) (NodeRef, error) {
	if len(path) == 0 {
		return ref, nil
	}
	node, ok := ref.node().(ContainerNode)
	if !ok {
		return nil, errors.New("must be container type")
	}
	child, err := node.getChild(path[0])
	if err != nil {
		return nil, err
	}
	newPath := path[1:]
	res, err := (*child).selectNodeByNodePath(newPath)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("unable to get node")
	}
	return res, err
}

func (ref *NodeReference) nodeType() NodeType {
	return ref.n.nodeType()
}
