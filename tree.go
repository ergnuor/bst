package bst

type Tree interface {
	GetRoot() Node
	Insert(...Payload)
	Delete(...Payload)
	Search(Payload) Payload
}
