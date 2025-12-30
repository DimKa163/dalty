package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraph(t *testing.T) {
	graph := NewGraph()
	nodeA := &Node{
		ID:    "A",
		Value: "A",
	}
	nodeB := &Node{
		ID:    "B",
		Value: "B",
	}
	nodeC := &Node{
		ID:    "C",
		Value: "C",
	}
	nodeD := &Node{
		ID:    "D",
		Value: "D",
	}
	nodeP := &Node{
		ID:     "P",
		Value:  "P",
		Master: true,
	}
	graph.AddNode(nodeA)
	graph.AddNode(nodeB)
	graph.AddNode(nodeC)
	graph.AddNode(nodeD)
	graph.AddNode(nodeP)
	graph.AddEdge(nodeP, nodeA, 0)
	graph.AddEdge(nodeP, nodeB, 0)
	graph.AddEdge(nodeP, nodeC, 0)
	graph.AddEdge(nodeP, nodeD, 0)

	n, ok := graph.Find(nodeA.ID)

	assert.True(t, ok)
	assert.Equal(t, nodeA, n)

	nodes := graph.AllIncomeTo(n)

	assert.Equal(t, 1, len(nodes))

	n, ok = graph.Find(nodeP.ID)

	assert.True(t, ok)
	assert.Equal(t, nodeP, n)

	nodes = graph.AllOutcomeFrom(n)
	assert.Equal(t, 4, len(nodes))
}
