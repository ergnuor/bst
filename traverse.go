package bst

import "math"

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
	n := t.Root()
	nodeStack := make([]Node, 0, t.MaxHeight())

	for l := len(nodeStack); n != nil || l > 0; l = len(nodeStack) {
		if n != nil {
			v.Visit(n.Payload())
			nodeStack = append(nodeStack, n)
			n = n.Left()
			continue
		}

		n = nodeStack[l-1]
		n = n.Right()
		nodeStack = nodeStack[:l-1]
	}
}

func TraverseInOrder(t Tree, v Visitor) {
	n := t.Root()
	nodeStack := make([]Node, 0, t.MaxHeight())

	for l := len(nodeStack); n != nil || l > 0; l = len(nodeStack) {
		if n != nil {
			nodeStack = append(nodeStack, n)
			n = n.Left()
			continue
		}

		n = nodeStack[l-1]
		v.Visit(n.Payload())
		n = n.Right()
		nodeStack = nodeStack[:l-1]
	}
}

func TraversePostOrder(t Tree, v Visitor) {
	n := t.Root()
	nodeStack := make([]Node, 0, t.MaxHeight())
	var lastVisitedNode Node

	for l := len(nodeStack); n != nil || l > 0; l = len(nodeStack) {
		if n != nil {
			nodeStack = append(nodeStack, n)
			n = n.Left()
			continue
		}

		n = nodeStack[l-1]
		if n.Right() != nil && n.Right() != lastVisitedNode {
			n = n.Right()
		} else {
			v.Visit(n.Payload())
			lastVisitedNode = n
			nodeStack = nodeStack[:l-1]
			n = nil
		}
	}
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
