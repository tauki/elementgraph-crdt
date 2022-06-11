package graph

import "github.com/google/uuid"

func (g *T) AddNode(node *Node) bool {
	if g.NodeExists(node) {
		return false
	}

	g.List[node.ID] = node
	return true
}

func (g *T) RemoveNode(node *Node) bool {
	if g.NodeExists(node) {
		for _, n := range g.List {
			for k, e := range n.Edges {
				if e.To.ID == node.ID {
					delete(n.Edges, k)
				}
			}
		}

		delete(g.List, node.ID)
		return true
	}

	return false
}

func (g *T) NodeExists(node *Node) bool {
	_, ok := g.List[node.ID]
	return ok
}

func (g *T) GetNode(id uuid.UUID) *Node {
	node, _ := g.List[id]
	return node
}

func (g *T) AddEdge(edge *Edge) bool {
	if g.EdgeExists(edge) ||
		!g.NodeExists(edge.To) ||
		!g.NodeExists(edge.From) {
		return false
	}
	g.List[edge.From.ID].Edges[edge.ID] = edge
	return true
}

func (g *T) RemoveEdge(edge *Edge) bool {
	if g.NodeExists(edge.From) {
		if g.EdgeExists(edge) {
			delete(g.List[edge.From.ID].Edges, edge.ID)
			return true
		}
	}

	return false
}

func (g *T) EdgeExists(edge *Edge) bool {
	if node, ok := g.List[edge.From.ID]; ok {
		_, ok = node.Edges[edge.ID]
		return ok
	}
	return false
}
