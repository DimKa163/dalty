package graph

type Graph struct {
	nodes map[string]*Node
	EdgeList
}

func NewGraph() *Graph {
	return &Graph{
		nodes:    make(map[string]*Node),
		EdgeList: make(EdgeList),
	}
}

func (g *Graph) Find(nodeID string) (*Node, bool) {
	n, ok := g.nodes[nodeID]
	return n, ok
}

func (g *Graph) AddNode(n *Node) {
	_, ok := g.nodes[n.ID]
	if ok {
		return
	}
	g.nodes[n.ID] = n
}
