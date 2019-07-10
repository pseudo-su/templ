package tree

import (
	"errors"
)

// func mergeAllNodes(trees ...NodeRef) (NodeRef, error) {
// 	return nil, nil
// }

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// nolint:gocyclo
func mergeObjectNodes(t1, t2 NodeRef) (NodeRef, error) {
	n1 := t1.node().(*ObjectNode)
	n2 := t2.node().(*ObjectNode)

	allKeys := []string{}
	allKeys = append(allKeys, n1.sortedKeys...)
	for _, str := range n2.sortedKeys {
		allKeys = appendIfMissing(allKeys, str)
	}
	allValues := map[string]NodeRef{}
	for _, key := range allKeys {
		baseVal, hasBaseVal := n1.raw[key]
		overlayVal, hasOverlayVal := n2.raw[key]
		switch {
		case hasBaseVal && baseVal != nil && hasOverlayVal && overlayVal != nil:
			nodeRef, err := mergeNodes(baseVal, overlayVal)
			if err != nil {
				return nil, err
			}
			allValues[key] = nodeRef
		case hasBaseVal && baseVal != nil:
			allValues[key] = baseVal
		case hasOverlayVal && overlayVal != nil:
			allValues[key] = overlayVal
		default:
			return nil, errors.New("unknown error merging object nodes")
		}
	}

	t1.swapNode(&ObjectNode{raw: allValues, sortedKeys: allKeys})
	return t1, nil
}

func mergeArrayNodes(t1, t2 NodeRef) (NodeRef, error) {
	n1 := t1.node().(*ArrayNode)
	n2 := t2.node().(*ArrayNode)

	allItems := []NodeRef{}
	allItems = append(allItems, n1.raw...)
	allItems = append(allItems, n2.raw...)

	t1.swapNode(&ArrayNode{raw: allItems})
	return t1, nil
}

func MergeTrees(t1, t2 NodeRef) (NodeRef, error) {
	return mergeNodes(t1, t2)
}

func mergeNodes(baseTree, overlayTree NodeRef) (NodeRef, error) {
	t1Type := baseTree.nodeType()
	t2Type := overlayTree.nodeType()
	switch {
	case t1Type == TObjectNode || t2Type == TObjectNode:
		if t1Type != t2Type {
			return nil, errors.New("mismatched node types when merging trees")
		}
		return mergeObjectNodes(baseTree, overlayTree)
	case t1Type == TArrayNode || t2Type == TArrayNode:
		if t1Type != t2Type {
			return nil, errors.New("mismatched node types when merging trees")
		}
		return mergeArrayNodes(baseTree, overlayTree)
	case t1Type == TNullNode || t2Type == TNullNode:
		// TODO: at least one node is null to take the other one
		return overlayTree, nil
	default:
		return overlayTree, nil

		// default:
		// 	return nil, errors.New("unknown node type")
	}
}
