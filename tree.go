package bst

type Tree interface {
	Root() Node
	Insert(...Payload)
	Delete(...Payload)
	Search(Payload) Payload
}
