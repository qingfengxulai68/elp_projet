package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
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
	
	fmt.Println("Avant Louvain :")
	graph.PrintGraph()

	// Temps d'exécution
	start := time.Now()

	// Lancer l'algorithme de Louvain pour détecter les communautés
	graph.Louvain(100)

	// Afficher les communautés
	graph.DisplayCommunities()

	// Afficher le temps d'exécution
	elapsed := time.Since(start)
	fmt.Printf("Temps d'exécution: %s\n", elapsed)
}