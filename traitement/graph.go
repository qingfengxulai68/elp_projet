package traitement

import (
	"fmt"
)

// Graph représente un graphe non orienté avec des communautés
type Graph struct {
	AdjList   map[string][]string // Liste d'adjacence
	Communities map[string]int    // Mapping entre les nœuds et leurs communautés
}

// NewGraph crée un nouveau graphe
func NewGraph() *Graph {
	return &Graph{
		AdjList:   make(map[string][]string),
		Communities: make(map[string]int),
	}
}

// AddEdge ajoute une arête non orientée au graphe
func (g *Graph) AddEdge(node1, node2 string) {
	g.AdjList[node1] = append(g.AdjList[node1], node2)
	g.AdjList[node2] = append(g.AdjList[node2], node1)
}

// Print affiche la liste d'adjacence du graphe
func (g *Graph) Print() {
	for node, neighbors := range g.AdjList {
		fmt.Printf("%s: %v\n", node, neighbors)
	}
}