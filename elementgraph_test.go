package crdt

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tauki/crdt/graph"
	"github.com/tauki/crdt/twoPSet"
	"testing"
	"time"
)

func TestGetElementGraph(t *testing.T) {
	g := NewElementGraph()
	assert.NotNil(t, g)
	assert.NotNil(t, g.Graph)
	assert.NotNil(t, g.EdgeSet)
	assert.NotNil(t, g.NodeSet)
}

func TestElementGraph_AddNode(t *testing.T) {
	g := NewElementGraph()
	id := uuid.New()

	node := graph.NewNode(id, []byte("hello"))
	g.AddNode(node)

	assert.Equal(t, node, g.NodeSet.GetAddSet()[id].Payload)
}

func TestElementGraph_AddNode_Idempotency(t *testing.T) {
	g := NewElementGraph()

	id := uuid.New()
	node1 := graph.NewNode(id, []byte("hello"))
	node2 := graph.NewNode(id, []byte("world"))
	g.AddNode(node1)
	assert.Len(t, g.NodeSet.GetAddSet(), 1)
	g.AddNode(node2)
	assert.Len(t, g.NodeSet.GetAddSet(), 1)

	assert.Equal(t, node1, g.NodeSet.GetAddSet()[id].Payload)
}

func TestSBLEGraph_AddEdge(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))

	g.AddNode(node1)
	g.AddNode(node2)

	edge := graph.NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)

	assert.NotNil(t, g.EdgeSet.GetAddSet())
	assert.Equal(t, edge, g.EdgeSet.GetAddSet()[edge.ID].Payload)
}

func TestElementGraph_AddEdge_Idempotency(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))

	g.AddNode(node1)
	g.AddNode(node2)

	edgeID := uuid.New()
	edge1 := graph.NewEdge(edgeID, node1, node2)
	edge2 := graph.NewEdge(edgeID, node1, node2)

	g.AddEdge(edge1)
	assert.Len(t, g.EdgeSet.GetAddSet(), 1)

	g.AddEdge(edge2)
	assert.Len(t, g.EdgeSet.GetAddSet(), 1)

	assert.Equal(t, edge1, g.EdgeSet.GetAddSet()[edgeID].Payload)
}

func TestElementGraph_RemoveNode(t *testing.T) {
	g := NewElementGraph()

	node := graph.NewNode(uuid.New(), []byte("hello"))
	g.AddNode(node)
	assert.NotNil(t, g.NodeSet.GetAddSet())

	g.RemoveNode(node)
	assert.NotNil(t, g.NodeSet.GetRemoveSet())

	assert.Equal(t, node, g.NodeSet.GetRemoveSet()[node.ID].Payload)
}

func TestElementGraph_RemoveNode_WithEdges(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	g.AddNode(node1)
	g.AddNode(node2)
	assert.NotNil(t, g.NodeSet.GetAddSet())

	edge := graph.NewEdge(uuid.New(), node2, node1)
	g.AddEdge(edge)

	g.RemoveNode(node1)
	assert.NotNil(t, g.NodeSet.GetRemoveSet())

	assert.Equal(t, node1, g.NodeSet.GetRemoveSet()[node1.ID].Payload)
	assert.False(t, g.Graph.EdgeExists(edge))
}

func TestElementGraph_RemoveNode_RemoveFromNodeSetFail(t *testing.T) {
	g := NewElementGraph()
	mockSet := &twoPSet.MockTwoPSet{}
	g.NodeSet = mockSet

	mockSet.On("Remove",
		mock.AnythingOfType("uuid.UUID"),
	).Return(errors.New("error"))

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	edge := graph.NewEdge(uuid.New(), node1, node2)

	g.Graph.AddNode(node1)
	g.Graph.AddNode(node2)
	g.Graph.AddEdge(edge)

	g.RemoveNode(node1)

	assert.True(t, g.Graph.NodeExists(node1))
}

func TestElementGraph_RemoveNode_Idempotency(t *testing.T) {
	g := NewElementGraph()

	nodeID := uuid.New()
	node1 := graph.NewNode(nodeID, []byte("hello"))
	g.AddNode(node1)
	assert.NotNil(t, g.NodeSet.GetAddSet())

	node2 := graph.NewNode(nodeID, []byte("world"))

	g.RemoveNode(node1)
	g.RemoveNode(node2)
	assert.NotNil(t, g.NodeSet.GetRemoveSet())

	assert.Equal(t, node1, g.NodeSet.GetRemoveSet()[nodeID].Payload)
	assert.Len(t, g.NodeSet.GetRemoveSet(), 1)
}

func TestElementGraph_RemoveEdge(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))

	g.AddNode(node1)
	g.AddNode(node2)

	edge := graph.NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)
	assert.NotNil(t, g.EdgeSet.GetAddSet())

	g.RemoveEdge(edge)
	assert.NotNil(t, g.EdgeSet.GetRemoveSet())
	assert.Equal(t, edge, g.EdgeSet.GetRemoveSet()[edge.ID].Payload)
}

func TestElementGraph_RemoveEdge_Idempotency(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))

	g.AddNode(node1)
	g.AddNode(node2)

	edgeID := uuid.New()
	edge1 := graph.NewEdge(edgeID, node1, node2)
	g.AddEdge(edge1)
	assert.NotNil(t, g.EdgeSet.GetAddSet())

	edge2 := graph.NewEdge(edgeID, node1, node2)
	g.RemoveEdge(edge1)
	g.RemoveEdge(edge2)

	assert.NotNil(t, g.EdgeSet.GetRemoveSet())
	assert.Equal(t, edge1, g.EdgeSet.GetRemoveSet()[edgeID].Payload)
	assert.Len(t, g.EdgeSet.GetRemoveSet(), 1)
}

func TestElementGraph_RemoveEdge_RemoveFromEdgeSetFail(t *testing.T) {
	g := NewElementGraph()
	mockSet := &twoPSet.MockTwoPSet{}
	g.EdgeSet = mockSet

	mockSet.On("Remove",
		mock.AnythingOfType("uuid.UUID"),
	).Return(errors.New("error"))

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	edge := graph.NewEdge(uuid.New(), node1, node2)

	g.Graph.AddNode(node1)
	g.Graph.AddNode(node2)
	g.Graph.AddEdge(edge)

	g.RemoveEdge(edge)

	assert.True(t, g.Graph.EdgeExists(edge))
}

func TestElementGraph_Merge(t *testing.T) {
	g1 := NewElementGraph()
	g2 := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	edge := graph.NewEdge(uuid.New(), node1, node2)

	g1.AddNode(node1)
	g1.AddNode(node2)
	g1.AddEdge(edge)

	g2.Merge(g1)
	assert.True(t, g2.Graph.NodeExists(node1))
	assert.True(t, g2.Graph.NodeExists(node2))
	assert.True(t, g2.Graph.EdgeExists(edge))
	assert.Equal(t, node1, g2.NodeSet.GetAddSet()[node1.ID].Payload)
	assert.Equal(t, node2, g2.NodeSet.GetAddSet()[node2.ID].Payload)
	assert.Equal(t, edge, g2.EdgeSet.GetAddSet()[edge.ID].Payload)
}

func TestElementGraph_Merge_Conflict(t *testing.T) {
	g1 := NewElementGraph()
	g2 := NewElementGraph()

	nodeID := uuid.New()
	edgeID := uuid.New()
	node1 := graph.NewNode(nodeID, []byte("hello"))
	node2 := graph.NewNode(nodeID, []byte("world"))
	edge1 := graph.NewEdge(edgeID, node1, node1)
	edge2 := graph.NewEdge(edgeID, node2, node2)

	g1.AddNode(node1)
	g1.AddEdge(edge1)

	time.Sleep(1)
	g2.AddNode(node2)
	g2.AddEdge(edge2)

	g1.Merge(g2)
	assert.True(t, g1.Graph.NodeExists(node2))
	assert.True(t, g1.Graph.EdgeExists(edge2))
	assert.Equal(t, node2, g1.NodeSet.GetAddSet()[nodeID].Payload)
	assert.Equal(t, edge2, g1.EdgeSet.GetAddSet()[edgeID].Payload)
}

func TestElementGraph_Merge_RemoveNodeConflictResolvedByTimestamp(t *testing.T) {
	g1 := NewElementGraph()
	g2 := NewElementGraph()

	nodeID := uuid.New()
	edgeID := uuid.New()
	node1 := graph.NewNode(nodeID, []byte("hello"))
	node2 := graph.NewNode(nodeID, []byte("world"))
	edge1 := graph.NewEdge(edgeID, node1, node1)
	edge2 := graph.NewEdge(edgeID, node2, node2)

	g1.AddNode(node1)
	g1.AddEdge(edge1)

	time.Sleep(1)
	g2.AddNode(node2)
	g2.AddEdge(edge2)
	time.Sleep(1)
	g2.RemoveNode(node2)

	g1.Merge(g2)
	assert.False(t, g1.Graph.NodeExists(node2))
	assert.False(t, g1.Graph.EdgeExists(edge2))
}

func TestElementGraph_Merge_RemoveEdgeConflictResolvedByTimestamp(t *testing.T) {
	g1 := NewElementGraph()
	g2 := NewElementGraph()

	nodeID := uuid.New()
	edgeID := uuid.New()
	node1 := graph.NewNode(nodeID, []byte("hello"))
	node2 := graph.NewNode(nodeID, []byte("world"))
	edge1 := graph.NewEdge(edgeID, node1, node1)
	edge2 := graph.NewEdge(edgeID, node2, node2)

	g1.AddNode(node1)
	g1.AddEdge(edge1)

	time.Sleep(1)
	g2.AddNode(node2)
	g2.AddEdge(edge2)
	time.Sleep(1)
	g2.RemoveEdge(edge2)

	g1.Merge(g2)
	assert.True(t, g1.Graph.NodeExists(node2))
	assert.False(t, g1.Graph.EdgeExists(edge2))
}

func TestElementGraph_RegenerateGraph(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	g.AddNode(node1)
	g.AddNode(node2)

	edge := graph.NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)

	g.Graph = graph.New()
	assert.False(t, g.Graph.NodeExists(node1))
	assert.False(t, g.Graph.NodeExists(node2))
	assert.False(t, g.Graph.EdgeExists(edge))

	g.RegenerateGraph()
	assert.True(t, g.Graph.NodeExists(node1))
	assert.True(t, g.Graph.NodeExists(node2))
	assert.True(t, g.Graph.EdgeExists(edge))
}

func TestElementGraph_RegenerateGraph_RemovedNodeSync(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	g.AddNode(node1)
	g.AddNode(node2)

	time.Sleep(1)
	g.RemoveNode(node1)

	g.Graph = graph.New()
	assert.False(t, g.Graph.NodeExists(node1))
	assert.False(t, g.Graph.NodeExists(node2))

	g.RegenerateGraph()
	assert.False(t, g.Graph.NodeExists(node1))
	assert.True(t, g.Graph.NodeExists(node2))
}

func TestElementGraph_RegenerateGraph_RemovedEdgeSync(t *testing.T) {
	g := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte("hello"))
	node2 := graph.NewNode(uuid.New(), []byte("world"))
	g.AddNode(node1)
	g.AddNode(node2)

	edge := graph.NewEdge(uuid.New(), node1, node2)
	g.AddEdge(edge)

	time.Sleep(1)
	g.RemoveEdge(edge)

	g.Graph = graph.New()
	assert.False(t, g.Graph.NodeExists(node1))
	assert.False(t, g.Graph.NodeExists(node2))
	assert.False(t, g.Graph.EdgeExists(edge))

	g.RegenerateGraph()
	assert.True(t, g.Graph.NodeExists(node1))
	assert.True(t, g.Graph.NodeExists(node2))
	assert.False(t, g.Graph.EdgeExists(edge))
}
