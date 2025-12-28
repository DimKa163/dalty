package core

import (
	"container/list"
	"context"
	"github.com/DimKa163/dalty/internal/shared"
	"github.com/beevik/guid"
	"sync"
)

type Graph struct {
	nodes    map[guid.Guid]*Node
	list     []*Node
	edgeList EdgeList
}

func NewGraph() *Graph {
	return &Graph{
		nodes:    make(map[guid.Guid]*Node),
		list:     make([]*Node, 0),
		edgeList: make(EdgeList),
	}
}

func (g *Graph) AddNode(n *Node) {
	_, ok := g.nodes[n.ID]
	if ok {
		return
	}
	g.nodes[n.ID] = n
	g.list = append(g.list, n)
}

func (g *Graph) AddEdge(from, to *Node, weight int) {
	g.edgeList.Add(from, to, weight)
}

func (g *Graph) Outgoing(n *Node) []*Node {
	edges := g.edgeList.OutcomeFrom(n)
	result := make([]*Node, len(edges))
	for i, e := range edges {
		result[i] = e.To
	}
	return result
}

func (g *Graph) Find(id *guid.Guid) (n *Node, ok bool) {
	n, ok = g.nodes[*id]
	return
}

func (g *Graph) Path(destination *Node) *Path {
	path := NewPath()
	queue := list.New()
	visited := make(map[guid.Guid]bool)
	subD := make(map[guid.Guid]int)
	subD[destination.ID] = 1
	queue.PushFront(&PathNode{
		Node: destination,
	})
	for queue.Len() > 0 {
		item := queue.Remove(queue.Front()).(*PathNode)
		_, ok := visited[item.ID]
		if ok {
			continue
		}
		visited[item.ID] = true
		item.Level = subD[item.ID]
		path.AddNode(item)
		// идём назад
		edges := g.edgeList.IncomeTo(item.Node)
		for _, edge := range edges {
			node := edge.From
			subD[node.ID] = item.Level + 1
			queue.PushFront(&PathNode{
				Level: item.Level,
				Node:  node,
				Next:  item,
			})
		}
	}
	return path
}

type NodeType int

const (
	WarehouseUNRECOGNIZED NodeType = iota
	WarehouseFREE
	WarehouseMAIN
	WarehouseCENTER
	WarehouseMALL
)

func (w NodeType) String() string {
	return []string{"UNRECOGNIZED", "FREE", "MAIN", "CENTER", "MALL"}[w]
}

func MapWarehouseType(code string) NodeType {
	switch code {
	case shared.WarehouseCategoryFREE:

		return WarehouseFREE
	case shared.WarehouseCategoryMAIN:

		return WarehouseMAIN
	case shared.WarehouseCategoryCENTRAL:
		return WarehouseCENTER
	case shared.WarehouseCategoryMALL:
		return WarehouseMALL
	default:
		return WarehouseUNRECOGNIZED
	}
}

type Node struct {
	ID                     guid.Guid
	Type                   NodeType
	Name                   string
	Code                   string
	AvailableRest          bool
	Address                string
	DescriptorGroup        string
	OnlyStockPickupAllowed bool
}

type Edge struct {
	From   *Node
	To     *Node
	Weight int
}

type EdgeList map[guid.Guid][]*Edge

func (el EdgeList) Add(from, to *Node, weight int) {
	edge := &Edge{from, to, weight}
	el.addIncomeAdd(edge)
	el.addOutcomeAdd(edge)
}

func (el EdgeList) IncomeTo(to *Node) []*Edge {
	result := make([]*Edge, 0)
	edges, ok := el[to.ID]
	if !ok {
		return result
	}
	for _, edge := range edges {
		if edge.To.ID != to.ID {
			continue
		}
		result = append(result, edge)
	}
	return result
}

func (el EdgeList) OutcomeFrom(from *Node) []*Edge {
	result := make([]*Edge, 0)
	edges, ok := el[from.ID]
	if !ok {
		return result
	}
	for _, edge := range edges {
		if edge.From.ID != from.ID {
			continue
		}
		result = append(result, edge)
	}
	return result
}

func (el EdgeList) Edges(n *Node) ([]*Edge, bool) {
	edges, ok := el[n.ID]
	if !ok {
		return nil, false
	}
	return edges, true
}
func (el EdgeList) addOutcomeAdd(edge *Edge) {
	_, ok := el[edge.From.ID]
	if !ok {
		el[edge.From.ID] = make([]*Edge, 0)
	}
	el[edge.From.ID] = append(el[edge.From.ID], edge)
}
func (el EdgeList) addIncomeAdd(edge *Edge) {
	_, ok := el[edge.To.ID]
	if !ok {
		el[edge.To.ID] = make([]*Edge, 0)
	}
	el[edge.To.ID] = append(el[edge.To.ID], edge)
}

type GraphContext struct {
	graph *Graph
	mutex *sync.RWMutex
}

func NewGraphContext() *GraphContext {
	return &GraphContext{
		mutex: &sync.RWMutex{},
	}
}

func (gc *GraphContext) Get(ctx context.Context) (*Graph, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		gc.mutex.RLock()
		defer gc.mutex.RUnlock()
		return gc.graph, nil
	}

}

func (gc *GraphContext) Update(graph *Graph) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.graph = graph
}
