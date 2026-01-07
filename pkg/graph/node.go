package graph

type Node struct {
	ID     string
	Value  any
	Master bool
}

func Cast[T any](n *Node) (T, bool) {
	v, ok := n.Value.(T)
	return v, ok
}
