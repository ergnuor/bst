package bst

type Node interface {
	GetPayload() Payload
	GetLeftChild() Node
	GetRightChild() Node
}
