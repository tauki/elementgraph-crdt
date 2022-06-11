package crdt

import (
	"github.com/tauki/crdt/graph"
	"github.com/tauki/crdt/twoPSet"
)

type ElementGraph struct {
	NodeSet twoPSet.TwoPSet
	EdgeSet twoPSet.TwoPSet
	Graph   graph.Graph
}

func NewElementGraph() *ElementGraph {
	return &ElementGraph{
		NodeSet: twoPSet.New(),
		EdgeSet: twoPSet.New(),
		Graph:   graph.New(),
	}
}

func (s *ElementGraph) AddNode(node *graph.Node) {
	if s.Graph.AddNode(node) {
		s.NodeSet.Add(node.ID, node)
	}
}

func (s *ElementGraph) AddEdge(edge *graph.Edge) {
	if s.Graph.AddEdge(edge) {
		s.EdgeSet.Add(edge.ID, edge)
	}
}

func (s *ElementGraph) RemoveNode(node *graph.Node) {
	if s.Graph.RemoveNode(node) {
		if err := s.NodeSet.Remove(node.ID); err != nil {
			s.Graph.AddNode(node)
			for _, v := range node.Edges {
				s.Graph.AddEdge(v)
			}
		}
	}
}

func (s *ElementGraph) RemoveEdge(edge *graph.Edge) {
	if s.Graph.RemoveEdge(edge) {
		if err := s.EdgeSet.Remove(edge.ID); err != nil {
			s.Graph.AddEdge(edge)
		}
	}
}

func (s *ElementGraph) Merge(g *ElementGraph) {
	s.NodeSet.Merge(g.NodeSet)
	s.EdgeSet.Merge(g.EdgeSet)
	s.RegenerateGraph()
}

func (s *ElementGraph) RegenerateGraph() {
	s.Graph = graph.New()

	addSet := s.NodeSet.GetAddSet()
	removeSet := s.NodeSet.GetRemoveSet()

	for k, v := range addSet {
		removedNode, ok := removeSet[k]
		if ok {
			if removedNode.Timestamp.After(v.Timestamp) {
				continue
			}
		}

		node := v.Payload.(*graph.Node)
		s.Graph.AddNode(graph.NewNode(node.ID, node.Payload))
	}

	addSet = s.EdgeSet.GetAddSet()
	removeSet = s.EdgeSet.GetRemoveSet()

	for k, v := range addSet {
		removedEdge, ok := removeSet[k]
		if ok {
			if removedEdge.Timestamp.After(v.Timestamp) {
				continue
			}
		}
		s.Graph.AddEdge(v.Payload.(*graph.Edge))
	}
}
