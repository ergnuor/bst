package redblack

import (
	"github.com/ergnuor/bst"
)

type node struct {
	payload bst.Payload
	left    *node
	right   *node
	blk     bool
}

func (n *node) Payload() bst.Payload {
	return n.payload
}

func (n *node) Left() bst.Node {
	if n.left == nil {
		return nil
	}

	return n.left
}

func (n *node) Right() bst.Node {
	if n.right == nil {
		return nil
	}

	return n.right
}
