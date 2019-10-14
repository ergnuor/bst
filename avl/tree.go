package avl

import (
	"math"

	"github.com/ergnuor/bst"
)

type tree struct {
	root *node
	cnt  int
}

func (t *tree) Root() bst.Node {
	if t.root == nil {
		return nil
	}

	return t.root
}

func (t *tree) Insert(pls ...bst.Payload) {
	for _, pl := range pls {
		n, path := t.findNode(pl)
		if *n != nil {
			continue
		}

		*n = t.newNode(pl)
		path = append(path, *n)
		t.fixInsertion(path)

		t.cnt++
	}
}

func (t *tree) Delete(pls ...bst.Payload) {
	for _, pl := range pls {
		n, path := t.findNode(pl)
		if *n == nil {
			continue
		}

		path = t.pickInOrderSuccessor(path)
		(*n).payload = path[len(path)-1].payload

		target := path[len(path)-1]
		linkToTarget := t.getLinkToRemovableNode(path)

		t.fixDeletion(path)

		t.doDeleteNode(target, linkToTarget)

		t.cnt--
	}
}

func (t *tree) Search(payload bst.Payload) bst.Payload {
	n, _ := t.findNode(payload)

	if *n != nil {
		return (*n).payload
	}

	return nil
}

func (t *tree) Count() int {
	return t.cnt
}

func (t *tree) MaxHeight() int {
	return int(math.Ceil(math.Log2(float64(t.cnt + 1))))
}

func New() bst.Tree {
	return &tree{}
}

func (t *tree) fixInsertion(path []*node) {
	var pivotParent *node

	for i := len(path) - 2; i >= 0; i-- {
		pivotParent = nil
		if i != 0 {
			pivotParent = path[i-1]
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
				t.fixLeftHeavy(path[i], path[i+1], path[i+2], pivotParent)
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
				t.fixRightHeavy(path[i], path[i+1], path[i+2], pivotParent)
				break
			}
		}
	}
}

func (t *tree) fixDeletion(path []*node) {
	var pivotParent *node

	for i := len(path) - 2; i >= 0; i-- {
		pivotParent = nil

		if i != 0 {
			pivotParent = path[i-1]
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
				path[i] = t.fixRightHeavy(path[i], path[i].right, t.pickZ(path[i].right), pivotParent)
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
				path[i] = t.fixLeftHeavy(path[i], path[i].left, t.pickZ(path[i].left), pivotParent)
			}
		}
	}
}

func (*tree) pickZ(n *node) *node {
	if n.bf == 1 {
		return n.right
	}
	return n.left
}

func (t *tree) pickInOrderSuccessor(path []*node) []*node {
	i := len(path) - 1

	if path[i].right == nil {
		return path
	}

	path = append(path, path[i].right)
	i++

	for path[i].left != nil {
		path = append(path, path[i].left)
		i++
	}

	return path
}

func (t *tree) getLinkToRemovableNode(path []*node) **node {
	l := len(path) - 1

	if l == 0 {
		return &t.root
	}

	if path[l-1].left == path[l] {
		return &path[l-1].left
	}

	return &path[l-1].right
}

func (t *tree) doDeleteNode(target *node, linkToTarget **node) {
	if target.left != nil {
		*linkToTarget = target.left
	} else if target.right != nil {
		*linkToTarget = target.right
	} else {
		*linkToTarget = nil
	}
}

func (t *tree) findNode(pl bst.Payload) (**node, []*node) {
	h := t.MaxHeight() + 1
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

func (t *tree) fixLeftHeavy(x *node, y *node, z *node, pivotParent *node) *node {
	// left left
	if y.left == z {
		t.rotateRight(x, pivotParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	// left right
	t.rotateLeft(y, x)
	t.rotateRight(x, pivotParent)

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

func (t *tree) fixRightHeavy(x *node, y *node, z *node, pivotParent *node) *node {
	// Right right
	if y.right == z {
		t.rotateLeft(x, pivotParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	// Right left
	t.rotateRight(y, x)
	t.rotateLeft(x, pivotParent)

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

func (t *tree) rotateRight(x *node, pivotParent *node) {
	xLeftNode := x.left

	x.left = xLeftNode.right
	xLeftNode.right = x

	if pivotParent != nil {
		if pivotParent.right == x {
			pivotParent.right = xLeftNode
		} else {
			pivotParent.left = xLeftNode
		}
	}

	if x == t.root {
		t.root = xLeftNode
	}
}

func (t *tree) rotateLeft(x *node, pivotParent *node) {
	xRightNode := x.right

	x.right = xRightNode.left
	xRightNode.left = x

	if pivotParent != nil {
		if pivotParent.right == x {
			pivotParent.right = xRightNode
		} else {
			pivotParent.left = xRightNode
		}
	}

	if x == t.root {
		t.root = xRightNode
	}
}

func (t *tree) newNode(pl bst.Payload) *node {
	return &node{
		pl,
		nil,
		nil,
		0,
	}
}
