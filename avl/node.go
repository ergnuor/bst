package avl

import (
	"github.com/ergnuor/bst"
)

type node struct {
	payload bst.Payload
	left    *node
	right   *node
	bf      int8
}

func (n *node) GetPayload() bst.Payload {
	return n.payload
}

func (n *node) GetLeftChild() bst.Node {
	if n.left == nil {
		return nil
	}

	return n.left
}

func (n *node) GetRightChild() bst.Node {
	if n.right == nil {
		return nil
	}

	return n.right
}
