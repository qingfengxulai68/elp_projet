package main

import (
	"fmt"
	// "os"
	// "bufio"
	// "strings"
	"elp_projet/traitement"
)

// func main() {
// 	// Ouvrir le fichier
// 	file, err := os.Open("amazon_ungraph.txt")
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Initialiser le graphe
// 	graph := traitement.Graph{AdjList: make(map[string][]string)}

// 	// Lire le fichier ligne par ligne
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		nodes := strings.Fields(line) // Découper la ligne par espace
// 		if len(nodes) == 2 {
// 			graph.AddEdge(nodes[0], nodes[1]) // Ajouter l'arête
// 		}
// 	}
// 	graph.Print()

// 	// Vérifier les erreurs éventuelles
// 	if err := scanner.Err(); err != nil {
// 		fmt.Println("Erreur lors de la lecture :", err)
// 	}
// }

// package main

// import (
// 	"fmt"
// 	"louvain/traitement"
// )

func main() {

	// Créer un nouveau graphe
	graph := traitement.NewGraph()

	// Ajouter des arêtes
	graph.AddEdge("A", "B")
	graph.AddEdge("A", "C")
	graph.AddEdge("B", "C")
	graph.AddEdge("D", "E")
	graph.AddEdge("E", "F")
	graph.AddEdge("F", "D")

	// Afficher le graphe avant Louvain
	fmt.Println("Graph avant Louvain:")
	graph.Print()

	// Appliquer l'algorithme Louvain pour détecter les communautés
	graph.Louvain()

	// Afficher les communautés après l'algorithme Louvain
	fmt.Println("\nCommunautés après Louvain:")
	for node, community := range graph.Communities {
		fmt.Printf("Node %s est dans la communauté %d\n", node, community)
	}
}
