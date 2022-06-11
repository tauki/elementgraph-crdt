package graph

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestT_AddNode(t *testing.T) {
	g := New()
	id := uuid.New()
	node := NewNode(id, []byte{})

	ok := g.AddNode(node)
	assert.True(t, ok)
	assert.True(t, g.NodeExists(node))
	assert.Equal(t, node, g.List[node.ID])
}

func TestT_AddNode_AlreadyExists(t *testing.T) {
	g := New()
	node := NewNode(uuid.New(), []byte{})

	assert.True(t, g.AddNode(node))
	assert.False(t, g.AddNode(node))
}

func TestT_RemoveNode(t *testing.T) {
	g := New()
	node := NewNode(uuid.New(), []byte{})

	g.AddNode(node)

	assert.True(t, g.RemoveNode(node))
	assert.False(t, g.NodeExists(node))
}

func TestT_RemoveNode_NodeDoesntExist(t *testing.T) {
	g := New()
	node := NewNode(uuid.New(), []byte{})

	assert.False(t, g.RemoveNode(node))
}

func TestT_RemoveNode_PointedEdgesRemoved(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)

	assert.True(t, g.RemoveNode(node2))
	assert.False(t, g.EdgeExists(edge))
}

func TestT_NodeExists(t *testing.T) {
	g := New()
	node := NewNode(uuid.New(), []byte{})
	g.AddNode(node)
	assert.True(t, g.NodeExists(node))
}

func TestT_NodeExists_NotFound(t *testing.T) {
	g := New()
	node := NewNode(uuid.New(), []byte{})
	assert.False(t, g.NodeExists(node))
}

func TestT_GetNode(t *testing.T) {
	g := New()
	id := uuid.New()
	node := NewNode(id, []byte{})
	g.AddNode(node)

	node = g.GetNode(id)
	assert.NotNil(t, node)
	assert.Equal(t, id, node.ID)
}

func TestT_GetNode_NotFound(t *testing.T) {
	g := New()

	node := g.GetNode(uuid.New())
	assert.Nil(t, node)
}

func TestT_AddEdge(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	assert.True(t, g.AddEdge(edge))

	assert.True(t, g.EdgeExists(edge))
	assert.Equal(t, edge, g.List[edge.From.ID].Edges[edge.ID])
}

func TestT_AddEdge_AlreadyExists(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	assert.True(t, g.AddEdge(edge))
	assert.False(t, g.AddEdge(edge))
}

func TestT_AddEdge_ToNodeDoesntExist(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)

	edge := NewEdge(uuid.New(), node1, node2)
	assert.False(t, g.AddEdge(edge))
}

func TestT_AddEdge_FromNodeDoesntExist(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	assert.False(t, g.AddEdge(edge))
}

func TestT_RemoveEdge(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)
	assert.True(t, g.EdgeExists(edge))

	assert.True(t, g.RemoveEdge(edge))
	assert.False(t, g.EdgeExists(edge))

	assert.Empty(t, g.List[edge.From.ID].Edges)
}

func TestT_RemoveEdge_NodeDoesntExist(t *testing.T) {
	g := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	edge := NewEdge(uuid.New(), node1, node2)

	assert.False(t, g.RemoveEdge(edge))
	assert.False(t, g.NodeExists(node1))
}

func TestT_RemoveEdge_EdgeDoesntExist(t *testing.T) {
	g := New()

	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})
	edge := NewEdge(uuid.New(), node1, node2)

	g.AddNode(node1)
	g.AddNode(node2)

	assert.False(t, g.RemoveEdge(edge))
	assert.True(t, g.NodeExists(node1))
	assert.True(t, g.NodeExists(node2))
}

func TestT_EdgeExists(t *testing.T) {
	g := New()
	node1 := NewNode(uuid.New(), []byte{})
	node2 := NewNode(uuid.New(), []byte{})

	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)

	assert.True(t, g.EdgeExists(edge))
}
