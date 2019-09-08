package bst

type Payload interface {
	Compare(Payload) int
}
