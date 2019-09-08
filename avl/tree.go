package avl

import (
	"math"

	"github.com/ergnuor/bst"
)

type AVL struct {
	root      *node
	nodeCount int
}

func (t AVL) GetRoot() bst.Node {
	return t.root
}

func (t *AVL) Insert(pls ...bst.Payload) {
	for _, pl := range pls {
		n, path := t.findNode(pl)
		if *n != nil {
			continue
		}

		var xParent *node
		*n = t.newNode(pl)
		path = append(path, *n)

		t.balanceInsertion(path, xParent)

		t.nodeCount++
	}
}

func (t *AVL) Delete(pls ...bst.Payload) {
	for _, pl := range pls {
		n, path := t.findNode(pl)
		if *n == nil {
			continue
		}

		n, path = t.pickInOrderSuccessor(n, path)

		target := path[len(path)-1]
		linkToTarget := t.getLinkToRemovableNode(path)

		t.balanceDeletion(path)

		t.doDeleteNode(target, linkToTarget)

		t.nodeCount--
	}
}

func (t *AVL) Search(payload bst.Payload) bst.Payload {
	n, _ := t.findNode(payload)

	if *n != nil {
		return (*n).payload
	}

	return nil
}

func New() bst.Tree {
	return &AVL{}
}

func (t *AVL) balanceInsertion(path []*node, xParent *node) {
	for i := len(path) - 2; i >= 0; i-- {
		xParent = nil
		if i != 0 {
			xParent = path[i-1]
		}

		if path[i].left == path[i+1] {
			if path[i].bf == 1 {
				path[i].bf = 0
				break
			}

			if path[i].bf == 0 {
				path[i].bf = -1
				continue
			}

			if path[i].bf == -1 {
				t.fixLeftHeavy(path[i], path[i+1], path[i+2], xParent)
				break
			}
		} else {
			if path[i].bf == -1 {
				path[i].bf = 0
				break
			}

			if path[i].bf == 0 {
				path[i].bf = 1
				continue
			}

			if path[i].bf == 1 {
				t.fixRightHeavy(path[i], path[i+1], path[i+2], xParent)
				break
			}
		}
	}
}

func (t *AVL) balanceDeletion(path []*node) {
	var x, y, z, xParent *node

	for i := len(path) - 2; i >= 0; i-- {
		x, y, z, xParent = nil, nil, nil, nil

		if i != 0 {
			xParent = path[i-1]
		}

		if path[i].left == path[i+1] {
			if path[i].bf == 0 {
				path[i].bf = 1
				break
			}

			if path[i].bf == -1 {
				path[i].bf = 0
				continue
			}

			if path[i].bf == 1 {
				x = path[i]
				y = path[i].right
			}
		} else {
			if path[i].bf == 0 {
				path[i].bf = -1
				break
			}

			if path[i].bf == 1 {
				path[i].bf = 0
				continue
			}

			if path[i].bf == -1 {
				x = path[i]
				y = path[i].left
			}
		}

		if x != nil && y != nil {
			if y.bf == 1 {
				z = y.right
			} else {
				z = y.left
			}

			if x.right == y {
				path[i] = t.fixRightHeavy(x, y, z, xParent)
			} else {
				path[i] = t.fixLeftHeavy(x, y, z, xParent)
			}
		}
	}
}

func (t *AVL) pickInOrderSuccessor(n **node, path []*node) (**node, []*node) {
	if (*n).left == nil || (*n).right == nil {
		return n, path
	}
	noteToSwapPayloadWith := *n

	n = &(*n).right
	path = append(path, *n)

	for (*n).left != nil {
		n = &(*n).left
		path = append(path, *n)
	}

	noteToSwapPayloadWith.payload = (*n).payload
	return n, path
}

func (t *AVL) getLinkToRemovableNode(path []*node) **node {
	i := len(path) - 1

	if len(path) == 1 {
		return &t.root
	}

	if path[i-1].left == path[i] {
		return &path[i-1].left
	}

	return &path[i-1].right
}

func (t *AVL) doDeleteNode(n *node, linkToNode **node) {
	if n.left != nil {
		*linkToNode = n.left
	} else if n.right != nil {
		*linkToNode = n.right
	} else {
		*linkToNode = nil
	}
}

func (t *AVL) findNode(pl bst.Payload) (**node, []*node) {
	h := int(math.Floor(math.Log2(float64(t.nodeCount+1)))) + 2
	path := make([]*node, 0, h)

	n := &t.root

	for *n != nil {
		path = append(path, *n)

		if (*n).payload.Compare(pl) == 0 {
			break
		}

		if (*n).payload.Compare(pl) == 1 {
			n = &(*n).left
		} else {
			n = &(*n).right
		}
	}

	return n, path
}

func (t *AVL) fixLeftHeavy(x *node, y *node, z *node, xParent *node) *node {
	// left left
	if y.left == z {
		t.rotateRight(x, xParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	// left right
	t.rotateLeft(y, x)
	t.rotateRight(x, xParent)

	y.bf = 0
	if z.bf == -1 {
		x.bf = 1
		z.bf = 0
	} else if z.bf == 1 {
		x.bf = 0
		z.bf = -1
	} else {
		x.bf = 0
		z.bf = 0
	}

	return x
}

func (t *AVL) fixRightHeavy(x *node, y *node, z *node, xParent *node) *node {
	// Right right
	if y.right == z {
		t.rotateLeft(x, xParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	// Right left
	t.rotateRight(y, x)
	t.rotateLeft(x, xParent)

	y.bf = 0
	if z.bf == -1 {
		x.bf = 1
		z.bf = 0
	} else if z.bf == 1 {
		x.bf = 0
		z.bf = -1
	} else {
		x.bf = 0
		z.bf = 0
	}

	return x
}

func (t *AVL) rotateRight(x *node, xParent *node) {
	xLeftNode := x.left

	x.left = xLeftNode.right
	xLeftNode.right = x

	if xParent != nil {
		if xParent.right == x {
			xParent.right = xLeftNode
		} else {
			xParent.left = xLeftNode
		}
	}

	if x == t.root {
		t.root = xLeftNode
	}
}

func (t *AVL) rotateLeft(x *node, xParent *node) {
	xRightNode := x.right

	x.right = xRightNode.left
	xRightNode.left = x

	if xParent != nil {
		if xParent.right == x {
			xParent.right = xRightNode
		} else {
			xParent.left = xRightNode
		}
	}

	if x == t.root {
		t.root = xRightNode
	}
}

func (t *AVL) newNode(pl bst.Payload) *node {
	return &node{
		pl,
		nil,
		nil,
		0,
	}
}
