package crdt

import (
	"github.com/google/uuid"
	"github.com/tauki/crdt/graph"
	"log"
)

func ExampleElementGraph_AddNode() {
	g := NewElementGraph()
	node := graph.NewNode(uuid.New(), []byte{})
	g.AddNode(node)

	log.Println(g.Graph.NodeExists(node))
}

func ExampleElementGraph_AddEdge() {
	g := NewElementGraph()
	node := graph.NewNode(uuid.New(), []byte{})
	g.AddNode(node)

	edge := graph.NewEdge(uuid.New(), node, node) // <-- self pointing edge
	g.AddEdge(edge)

	log.Println(g.Graph.EdgeExists(edge))
}

func ExampleElementGraph_RemoveNode() {
	g := NewElementGraph()
	node := graph.NewNode(uuid.New(), []byte{})
	g.AddNode(node)

	g.RemoveNode(node)
	log.Println(g.Graph.NodeExists(node))
}

func ExampleElementGraph_RemoveEdge() {
	g := NewElementGraph()
	node := graph.NewNode(uuid.New(), []byte{})
	g.AddNode(node)

	edge := graph.NewEdge(uuid.New(), node, node) // <-- self pointing edge
	g.AddEdge(edge)
	g.RemoveEdge(edge)

	log.Println(g.Graph.EdgeExists(edge))
}

func ExampleElementGraph_Merge() {
	g1 := NewElementGraph()
	g2 := NewElementGraph()

	node1 := graph.NewNode(uuid.New(), []byte{})
	node2 := graph.NewNode(uuid.New(), []byte{})

	g1.AddNode(node1)
	g2.AddNode(node2)

	g1.Merge(g2)

	log.Println(g1.Graph.NodeExists(node2))
}

func ExampleElementGraph_RegenerateGraph() {
	g := NewElementGraph()
	node := graph.NewNode(uuid.New(), []byte{})
	g.AddNode(node)

	log.Println(g.Graph.NodeExists(node))

	_ = g.NodeSet.Remove(node.ID)
	g.RegenerateGraph()

	log.Println(g.Graph.NodeExists(node))
}
