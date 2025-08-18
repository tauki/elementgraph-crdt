package graph

import (
	"github.com/google/uuid"
)

type Graph interface {
	AddNode(*Node) bool
	AddEdge(*Edge) bool
	RemoveNode(*Node) bool
	RemoveEdge(*Edge) bool
	GetNode(uuid.UUID) *Node
	NodeExists(*Node) bool
	EdgeExists(*Edge) bool
	FindPath(start *Node, end *Node) []*Node
}

type NodeSet map[uuid.UUID]*Node
type EdgeSet map[uuid.UUID]*Edge

type Node struct {
	ID      uuid.UUID
	Edges   EdgeSet
	Payload []byte
}

func (n Node) GetEdges() EdgeSet {
	return n.Edges
}

func NewNode(id uuid.UUID, payload []byte) *Node {
	return &Node{
		Edges:   make(EdgeSet),
		Payload: payload,
		ID:      id,
	}
}

type Edge struct {
	ID   uuid.UUID
	From *Node
	To   *Node
}

func (e Edge) Equal(edge *Edge) bool {
	return e.From == edge.From && e.To == edge.To
}

func NewEdge(id uuid.UUID, from, to *Node) *Edge {
	return &Edge{
		From: from,
		To:   to,
		ID:   id,
	}
}

type T struct {
	List NodeSet
}

func New() *T {
	return &T{
		List: make(NodeSet),
	}
}

var _ Graph = &T{}

func (g *T) FindPath(start *Node, end *Node) []*Node {

	path := make([]*Node, 0)

	if !g.NodeExists(start) || !g.NodeExists(end) {
		return path
	}

	path, _ = g.findPath(start, end, path)
	return path
}

func (g *T) findPath(start *Node, end *Node, path []*Node) ([]*Node, bool) {
	path = append(path, start)

	if start == end {
		return path, true
	}

	for _, edge := range start.Edges {
		newPath, ok := g.findPath(edge.To, end, path)
		if ok {
			return newPath, true
		}
	}
	return []*Node{}, false
}
