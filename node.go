package bst

type Node interface {
	Payload() Payload
	Left() Node
	Right() Node
}
