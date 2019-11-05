package redblack

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

		t.cnt--

		path = t.pickInOrderSuccessor(path)
		(*n).payload = path[len(path)-1].payload

		child := t.pickChild(path)
		path = t.replaceWithChild(path, child)

		if t.fixDeletionSimpleCases(path, child) {
			return
		}

		t.fixDeletionComplexCases(path, child)
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
	return 2 * int(math.Ceil(math.Log2(float64(t.cnt+1))))
}

func New() bst.Tree {
	return &tree{}
}

func (t *tree) fixDeletionSimpleCases(path []*node, child *node) bool {
	i := len(path) - 1

	if t.root == nil {
		return true
	}

	if t.root == child {
		child.blk = true
		return true
	}

	if path[i].blk == false {
		return true
	}

	if child != nil && child.blk == false {
		child.blk = true
		return true
	}

	return false
}

func (t *tree) fixDeletionComplexCases(path []*node, child *node) {
	var pivotParent *node

	i := len(path) - 1
	path[i] = child

case1:

	// CASE 1
	if path[i] == t.root {
		return
	}

	// CASE 2
	s := t.sibling(path, i)
	if s != nil && s.blk == false {
		path[i-1].blk = false
		s.blk = true

		pivotParent = nil
		if i-2 >= 0 {
			pivotParent = path[i-2]
		}
		if s == path[i-1].right {
			t.rotateLeft(path[i-1], pivotParent)
		} else {
			t.rotateRight(path[i-1], pivotParent)
		}

		i++
		if i == len(path) {
			path = append(path, path[i-1])
		} else {
			path[i] = path[i-1]
		}
		path[i-1] = path[i-2]
		path[i-2] = s
	}

	// CASE 3
	s = t.sibling(path, i)
	if path[i-1].blk == true && s != nil && s.blk == true && (s.left == nil || s.left.blk == true) && (s.right == nil || s.right.blk == true) {
		s.blk = false
		i--
		goto case1
	}

	// CASE 4
	if path[i-1].blk == false && s != nil && s.blk == true && (s.left == nil || s.left.blk == true) && (s.right == nil || s.right.blk == true) {
		path[i-1].blk = true
		s.blk = false
		return
	}

	// CASE 5
	if s != nil && s.blk == true {
		if path[i] == path[i-1].right && s.right != nil && s.right.blk == false && (s.left == nil || s.left.blk == true) {
			s.blk = false
			s.right.blk = true
			t.rotateLeft(s, path[i-1])
		} else if path[i] == path[i-1].left && s.left != nil && s.left.blk == false && (s.right == nil || s.right.blk == true) {
			s.blk = false
			s.left.blk = true
			t.rotateRight(s, path[i-1])
		}
	}

	// CASE 6
	s = t.sibling(path, i)
	if s != nil && s.blk == true {
		pivotParent = nil
		if i-2 >= 0 {
			pivotParent = path[i-2]
		}

		if path[i] == path[i-1].right && s.left != nil && s.left.blk == false {
			t.rotateRight(path[i-1], pivotParent)
			s.blk = path[i-1].blk
			path[i-1].blk = true
			s.left.blk = true
		} else if path[i] == path[i-1].left && s.right != nil && s.right.blk == false {
			t.rotateLeft(path[i-1], pivotParent)
			s.blk = path[i-1].blk
			path[i-1].blk = true
			s.right.blk = true
		}
	}
}

func (t *tree) pickChild(path []*node) *node {
	i := len(path) - 1

	if path[i].right != nil {
		return path[i].right
	}

	return path[i].left
}

func (t *tree) replaceWithChild(path []*node, child *node) []*node {
	i := len(path) - 1

	if path[i] == t.root {
		t.root = child
	} else if path[i] == path[i-1].right {
		path[i-1].right = child
	} else {
		path[i-1].left = child
	}

	return path
}

func (t *tree) sibling(path []*node, i int) *node {
	if i == 0 {
		return nil
	}

	if path[i] == path[i-1].right {
		return path[i-1].left
	}

	return path[i-1].right
}

func (t *tree) fixInsertion(path []*node) {
	var pivotParent *node

	for i := len(path) - 1; i >= 0; i-- {
		pivotParent = nil

		if path[i] == t.root {
			path[i].blk = true
			return
		}

		if path[i-1].blk == true {
			return
		}

		if path[i-2].left == path[i-1] {
			if path[i-2].right != nil && path[i-2].right.blk == false {
				path[i-2].blk = false
				path[i-2].right.blk = true
				path[i-1].blk = true
				i = i - 1
				continue
			}

			if path[i-1].right == path[i] {
				t.rotateLeft(path[i-1], path[i-2])
				path = t.swapPathNodes(path, i, i-1)
			}

			if i-3 >= 0 {
				pivotParent = path[i-3]
			}

			t.rotateRight(path[i-2], pivotParent)

			path[i-2].blk = false
			path[i-1].blk = true
			return
		} else {
			if path[i-2].left != nil && path[i-2].left.blk == false {
				path[i-2].blk = false
				path[i-2].left.blk = true
				path[i-1].blk = true
				i = i - 1
				continue
			}

			if path[i-1].left == path[i] {
				t.rotateRight(path[i-1], path[i-2])
				path = t.swapPathNodes(path, i, i-1)
			}

			if i-3 >= 0 {
				pivotParent = path[i-3]
			}
			t.rotateLeft(path[i-2], pivotParent)

			path[i-2].blk = false
			path[i-1].blk = true
			return
		}
	}
}

func (*tree) swapPathNodes(path []*node, i, j int) []*node {
	tmp := path[i]
	path[i] = path[j]
	path[j] = tmp
	return path
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
		false,
	}
}
