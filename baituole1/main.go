package main

import (
	"baituole1/traitement"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// type Graph struct {
// 	AdjList     map[int][]int
// 	Communities map[int]int
// }

// NewGraph crée un graphe vide
// func NewGraph() *Graph {
// 	return &Graph{
// 		AdjList:     make(map[int][]int),
// 		Communities: make(map[int]int),
// 	}
// }

func main() {
	// Créer un nouveau graphe
	graph := traitement.NewGraph()

	// Ouvrir le fichier
	file, err := os.Open("minigraph.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	// Lire le fichier ligne par ligne
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignorer les lignes vides ou les commentaires (lignes qui commencent par "#")
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		nodes := strings.Fields(line) // Découper la ligne par espace
		if len(nodes) == 2 {
			// Convertir les chaînes en entiers
			u, err := strconv.Atoi(nodes[0])
			if err != nil {
				fmt.Println("Erreur lors de la conversion de", nodes[0], "en entier:", err)
				continue
			}
			v, err := strconv.Atoi(nodes[1])
			if err != nil {
				fmt.Println("Erreur lors de la conversion de", nodes[1], "en entier:", err)
				continue
			}

			// Ajouter l'arête au graphe
			graph.AddEdge(u, v)
		}
	}

	// Vérifier si l'il y a eu des erreurs pendant la lecture
	if err := scanner.Err(); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
	}

	// graph := traitement.NewGraph()
	// graph.AddEdge(1, 2)
	// graph.AddEdge(1, 3)
	// graph.AddEdge(2, 3)
	// graph.AddEdge(4, 5)
	// graph.AddEdge(5, 6)
	// graph.AddEdge(6, 4)

	// fmt.Println("Avant Louvain :")
	// graph.PrintGraph()

	// 声明并初始化一个 sync.Mutex
	var mu sync.Mutex

	// Temps d'exécution
	start := time.Now()

	// Lancer l'algorithme de Louvain pour détecter les communautés
	graph.Louvain(10, &mu)

	fmt.Println("Communautés trouvées :")
	// Afficher les communautés
	graph.DisplayCommunities()

	// Afficher le temps d'exécution
	elapsed := time.Since(start)
	fmt.Printf("Temps d'exécution: %s\n", elapsed)
}
