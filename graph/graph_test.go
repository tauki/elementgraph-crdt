package graph

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	graph := New()
	assert.NotNil(t, graph)
	assert.IsType(t, NodeSet{}, graph.List)
	assert.Empty(t, graph.List)
}

func TestNewNode(t *testing.T) {
	id := uuid.New()
	payload := []byte("hello")
	node := NewNode(id, payload)
	assert.IsType(t, EdgeSet{}, node.Edges)
	assert.Empty(t, node.Edges)
	assert.Equal(t, id, node.ID)
	assert.Equal(t, payload, node.Payload)
}

func TestNewEdge(t *testing.T) {
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	id := uuid.New()
	edge := NewEdge(id, node1, node2)
	assert.Equal(t, edge.From, node1)
	assert.Equal(t, edge.To, node2)
}

func TestNode_GetEdges(t *testing.T) {
	id := uuid.New()
	node1 := NewNode(id, []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	edge := NewEdge(uuid.New(), node1, node2)

	node1.Edges[edge.ID] = edge

	edges := node1.GetEdges()

	assert.NotEmpty(t, edges)
	assert.Equal(t, edges[edge.ID], edge)
}

func TestEdge_Equal(t *testing.T) {
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	edge1 := NewEdge(uuid.New(), node1, node2)
	edge2 := NewEdge(uuid.New(), node1, node2)

	assert.True(t, edge1.Equal(edge2))
}

func TestT_FindPath(t *testing.T) {
	graph := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	node3 := NewNode(uuid.New(), []byte{})

	edge1 := NewEdge(uuid.New(), node1, node2)
	edge2 := NewEdge(uuid.New(), node2, node3)

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddNode(node3)

	graph.AddEdge(edge1)
	graph.AddEdge(edge2)

	path := graph.FindPath(node1, node3)

	assert.Equal(t, path[0].ID, node1.ID)
	assert.Equal(t, path[1].ID, node2.ID)
	assert.Equal(t, path[2].ID, node3.ID)
}

func TestT_FindPath_NoPath(t *testing.T) {
	graph := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	node3 := NewNode(uuid.New(), []byte{})

	edge1 := NewEdge(uuid.New(), node1, node2)

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddNode(node3)

	graph.AddEdge(edge1)

	path := graph.FindPath(node1, node3)
	assert.Empty(t, path)
}

func TestT_FindPath_StartNodeDoesntExist(t *testing.T) {
	graph := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	node3 := NewNode(uuid.New(), []byte{})

	edge1 := NewEdge(uuid.New(), node1, node2)

	graph.AddNode(node1)
	graph.AddNode(node2)

	graph.AddEdge(edge1)

	path := graph.FindPath(node3, node1)
	assert.Empty(t, path)
}

func TestT_FindPath_EndNodeDoesntExist(t *testing.T) {
	graph := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	node3 := NewNode(uuid.New(), []byte{})

	edge1 := NewEdge(uuid.New(), node1, node2)

	graph.AddNode(node1)
	graph.AddNode(node2)

	graph.AddEdge(edge1)

	path := graph.FindPath(node1, node3)
	assert.Empty(t, path)
}

func TestT_FindPathWithLoop(t *testing.T) {
	graph := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	node3 := NewNode(uuid.New(), []byte{})
	node4 := NewNode(uuid.New(), []byte{})

	edge12 := NewEdge(uuid.New(), node1, node2)
	edge23 := NewEdge(uuid.New(), node2, node3)
	edge32 := NewEdge(uuid.New(), node3, node2)
	edge14 := NewEdge(uuid.New(), node1, node4)

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddNode(node3)
	graph.AddNode(node4)

	graph.AddEdge(edge12)
	graph.AddEdge(edge23)
	graph.AddEdge(edge32)
	graph.AddEdge(edge14)

	path := graph.FindPath(node1, node4)

	assert.Equal(t, path[0].ID, node1.ID)
	assert.Equal(t, path[1].ID, node4.ID)
}
