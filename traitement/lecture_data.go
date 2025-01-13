package traitement

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Graph struct {
	adjList map[string][]string // Liste d'adjacence
}

func (g *Graph) AddEdge(node1, node2 string) {
	g.adjList[node1] = append(g.adjList[node1], node2)
	g.adjList[node2] = append(g.adjList[node2], node1) // Si le graphe est non orienté
}

// Afficher le graphe
func (g *Graph) Print() {
	for node, neighbors := range g.adjList {
		fmt.Printf("%s -> %v\n", node, neighbors)
	}
}


func main() {
	// Ouvrir le fichier
	file, err := os.Open("amazon_ungraph.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	// Initialiser le graphe
	graph := Graph{adjList: make(map[string][]string)}

	// Lire le fichier ligne par ligne
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nodes := strings.Fields(line) // Découper la ligne par espace
		if len(nodes) == 2 {
			graph.AddEdge(nodes[0], nodes[1]) // Ajouter l'arête
		}
	}
	// graph.Print()

	// Vérifier les erreurs éventuelles
	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture :", err)
	}
}
