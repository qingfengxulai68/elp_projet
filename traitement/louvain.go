package traitement

import "fmt"

// Modularity calcule la modularité d'une partition de graphes


func (g *Graph) Modularity() float64 {
	m := float64(0)
	for _, neighbors := range g.AdjList {
		m += float64(len(neighbors))
	}
	m /= 2

	Q := 0.0
	for node, neighbors := range g.AdjList {
		community := g.Communities[node]
		ki := float64(len(neighbors))
		for _, neighbor := range neighbors {
			kj := float64(len(g.AdjList[neighbor]))
			if g.Communities[neighbor] == community {
				Q += 1.0 - (ki * kj) / (2.0 * m)
			}
		}
	}
	res := (Q / (2 * m))
	// fmt.Printf("modularity : %f" , res )
	return res
}

// Louvain optimise la modularité et fusionne les communautés
func (g *Graph) Louvain() {
	// Initialiser les communautés (chaque nœud dans sa propre communauté)
	i := 0
	for node := range g.AdjList {
		g.Communities[node] = i
		fmt.Printf("Node %s : communauté %d \n", node, i)
		i += 1
	}

	// Répéter l'optimisation
	improvement := true
	for improvement {
		improvement = false
		// Optimiser la modularité localement
		for node := range g.AdjList {
			bestCommunity := g.Communities[node]
			bestModularity := g.Modularity()

			// Tester chaque communauté voisine
			for _, neighbor := range g.AdjList[node] {
				g.Communities[node] = g.Communities[neighbor]
				newModularity := g.Modularity()
				if newModularity > bestModularity {
					bestModularity = newModularity
					bestCommunity = g.Communities[neighbor]
					improvement = true
					fmt.Printf("Node %s : nouvelle communauté %d \n", node, i)
				}
			}

			// Assigner la meilleure communauté trouvée
			g.Communities[node] = bestCommunity
			fmt.Printf("Node %s : finale communauté %d \n", node, i)
		}

		// Fusionner les communautés
		if improvement {
			g.MergeCommunities()
		}
	}
}

// Fusionner les communautés en un nouveau graphe réduit
func (g *Graph) MergeCommunities() {
	// Créer un nouveau graphe pour les communautés fusionnées
	newGraph := NewGraph()

	// Parcourir chaque nœud et fusionner les communautés
	for node, community := range g.Communities {
		// Créer des arêtes entre les nœuds appartenant à la même communauté
		for _, neighbor := range g.AdjList[node] {
			// Si l'autre nœud est dans la même communauté, ajouter une arête
			if g.Communities[neighbor] == community {
				// Ajouter une arête entre les nœuds dans la même communauté
				newGraph.AddEdge(node, neighbor)
			}
		}
	}
	// Remplacer l'ancien graphe par le nouveau graphe fusionné
	*g = *newGraph
}

