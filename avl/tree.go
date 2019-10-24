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
	for i := len(path) - 2; i >= 0; i-- {
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
				t.fixLeftHeavyInsertion(path[i], t.pickParent(path, i))
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
				t.fixRightHeavyInsertion(path[i], t.pickParent(path, i))
				break
			}
		}
	}
}

func (t *tree) fixDeletion(path []*node) {
	for i := len(path) - 2; i >= 0; i-- {
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
				path[i] = t.fixRightHeavyDeletion(path[i], t.pickParent(path, i))

				if path[i].bf == -1 {
					break
				}
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
				path[i] = t.fixLeftHeavyDeletion(path[i], t.pickParent(path, i))

				if path[i].bf == 1 {
					break
				}
			}
		}
	}
}

func (*tree) pickParent(path []*node, i int) *node {
	if i != 0 {
		return path[i-1]
	}

	return nil
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

func (t *tree) fixLeftHeavyInsertion(x *node, xParent *node) {
	y := x.left

	// left left
	if y.bf == -1 {
		t.rotateRight(x, xParent)
		x.bf = 0
		y.bf = 0

		return
	}

	z := y.right

	// left right
	t.rotateLeft(y, x)
	t.rotateRight(x, xParent)

	t.fixZigZagBalanceFactor(z)
}

func (t *tree) fixRightHeavyInsertion(x *node, xParent *node) {
	y := x.right

	// Right right
	if y.bf == 1 {
		t.rotateLeft(x, xParent)
		x.bf = 0
		y.bf = 0

		return
	}

	z := y.left

	// Right left
	t.rotateRight(y, x)
	t.rotateLeft(x, xParent)

	t.fixZigZagBalanceFactor(z)
}

func (t *tree) fixLeftHeavyDeletion(x *node, xParent *node) *node {
	y := x.left

	if y.bf == -1 {
		t.rotateRight(x, xParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	if y.bf == 0 {
		t.rotateRight(x, xParent)
		y.bf = 1
		x.bf = -1

		return y
	}

	z := y.right

	// left right
	t.rotateLeft(y, x)
	t.rotateRight(x, xParent)

	t.fixZigZagBalanceFactor(z)

	return z
}

func (t *tree) fixRightHeavyDeletion(x *node, xParent *node) *node {
	y := x.right

	if y.bf == 1 {
		t.rotateLeft(x, xParent)
		x.bf = 0
		y.bf = 0

		return y
	}

	if y.bf == 0 {
		t.rotateLeft(x, xParent)
		y.bf = -1
		x.bf = 1

		return y
	}

	z := y.left

	// Right left
	t.rotateRight(y, x)
	t.rotateLeft(x, xParent)

	t.fixZigZagBalanceFactor(z)

	return z
}

func (*tree) fixZigZagBalanceFactor(p *node) {
	if p.bf == -1 {
		p.left.bf = 0
		p.right.bf = 1
	} else if p.bf == 1 {
		p.left.bf = -1
		p.right.bf = 0
	} else {
		p.left.bf = 0
		p.right.bf = 0
	}
	p.bf = 0
}

func (t *tree) rotateRight(x *node, xParent *node) {
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

func (t *tree) rotateLeft(x *node, xParent *node) {
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

func (t *tree) newNode(pl bst.Payload) *node {
	return &node{
		pl,
		nil,
		nil,
		0,
	}
}
