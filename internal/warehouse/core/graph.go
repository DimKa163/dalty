package core

import (
	"container/list"
	"context"
	"github.com/DimKa163/dalty/internal/graph"
)

type PathFinder struct {
	Context *graph.GraphContext
}

func NewPathFinder(graphContext *graph.GraphContext) *PathFinder {
	return &PathFinder{
		Context: graphContext,
	}
}

func (ws *PathFinder) Path(ctx context.Context, destination *graph.Node) (*Path, error) {
	path := NewPath()
	queue := list.New()
	visited := make(map[string]bool)
	subD := make(map[string]int)
	subD[destination.ID] = 1
	gr, err := ws.Context.Get(ctx)
	if err != nil {
		return nil, err
	}
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
		edges := gr.AllIncomeTo(item.Node)
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
	return path, nil
}
