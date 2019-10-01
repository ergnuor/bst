package bst

import (
	"math"
)

type order byte

const (
	pre order = iota
	in
	post
)

type Visitor interface {
	Visit(Payload)
}

type queue struct {
	f, b, l int
	q       []Node
}

func (q *queue) Enqueue(n Node) {
	if q.q[q.b] != nil || n == nil {
		return
	}

	q.q[q.b] = n
	q.b = int(math.Mod(float64(q.b+1), float64(q.l)))
}

func (q *queue) Dequeue() Node {
	if q.Empty() {
		return nil
	}

	n := q.q[q.f]
	q.q[q.f] = nil
	q.f = int(math.Mod(float64(q.f+1), float64(q.l)))

	return n
}

func (q *queue) Empty() bool {
	return q.q[q.f] == nil
}

func TraversePreOrder(t Tree, v Visitor) {
	orderTraverse(t, v, pre)
}

func TraverseInOrder(t Tree, v Visitor) {
	orderTraverse(t, v, in)
}

func TraversePostOrder(t Tree, v Visitor) {
	orderTraverse(t, v, post)
}

func TraverseBreadthFirst(t Tree, v Visitor) {
	l := int(math.Ceil(float64(t.Count()) / 2))
	q := &queue{l: l, q: make([]Node, l)}

	q.Enqueue(t.Root())

	for !q.Empty() {
		n := q.Dequeue()
		v.Visit(n.Payload())
		if n.Left() != nil {
			q.Enqueue(n.Left())
		}

		if n.Right() != nil {
			q.Enqueue(n.Right())
		}
	}
}

func orderTraverse(t Tree, v Visitor, o order) {
	n := t.Root()
	nodeStack := make([]Node, 0, t.MaxHeight())
	var lastPoppedNode Node

	for l := len(nodeStack); n != nil || l > 0; l = len(nodeStack) {
		if n != nil {
			if o&pre != 0 {
				v.Visit(n.Payload())
			}

			nodeStack = append(nodeStack, n)
			n = n.Left()
			continue
		}

		n = nodeStack[l-1]
		if o&in != 0 && (n.Right() == nil || n.Right() != lastPoppedNode) {
			v.Visit(n.Payload())
		}

		if n.Right() != nil && n.Right() != lastPoppedNode {
			n = n.Right()
			continue
		}

		if o&post != 0 {
			v.Visit(n.Payload())
		}
		lastPoppedNode = n
		nodeStack = nodeStack[:l-1]
		n = nil
	}
}
