package main

import (
	"bufio"
	"fmt"
	"gptandmain6/traitement"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 处理客户端请求
func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	// 读取客户端输入（格式："文件编号 Goroutine数量"）
	if !scanner.Scan() {
		fmt.Println("Erreur de lecture des paramètres")
		return
	}
	input := strings.Fields(scanner.Text())
	if len(input) != 2 {
		fmt.Println("Format invalide")
		return
	}

	// 解析文件编号
	fileNumber, err1 := strconv.Atoi(input[0])
	// 解析 Goroutine 数量
	goroutineCount, err2 := strconv.Atoi(input[1])

	if err1 != nil || err2 != nil || fileNumber < 1 || fileNumber > 3 || goroutineCount < 1 {
		fmt.Println("Paramètres invalides")
		return
	}

	// 生成文件名
	fileName := fmt.Sprintf("minigraph%d.txt", fileNumber)

	// 读取对应的文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(conn, "Erreur: impossible d'ouvrir %s\n", fileName)
		return
	}
	defer file.Close()

	// 创建一个新图
	graph := traitement.NewGraph()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
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

	if err := fileScanner.Err(); err != nil {
		fmt.Fprintf(conn, "Erreur lors de la lecture du fichier : %s\n", err)
		return
	}

	// 计算 Louvain 并发
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex // 互斥锁，防止并发写 map

	nodes := make([]int, 0, len(graph.AdjList))
	for node := range graph.AdjList {
		nodes = append(nodes, node)
	}

	// 切割任务
	batchSize := len(nodes) / goroutineCount
	if batchSize == 0 {
		batchSize = 1
	}

	// 启动 Goroutines
	for i := 0; i < goroutineCount; i++ {
		startIdx := i * batchSize
		endIdx := (i + 1) * batchSize
		if i == goroutineCount-1 {
			endIdx = len(nodes)
		}

		wg.Add(1)
		go func(nodesSubset []int) {
			defer wg.Done()
			mu.Lock() // 加锁，防止并发修改 map
			defer mu.Unlock()
			for _, node := range nodesSubset {
				currentCommunity := graph.Communities[node]
				bestCommunity := currentCommunity
				bestModularity := graph.Modularity()

				for _, neighbor := range graph.AdjList[node] {
					graph.Communities[node] = graph.Communities[neighbor]
					newModularity := graph.Modularity()
					if newModularity > bestModularity {
						bestModularity = newModularity
						bestCommunity = graph.Communities[neighbor]
					}
				}

				graph.Communities[node] = bestCommunity
			}
		}(nodes[startIdx:endIdx])
	}

	wg.Wait()

	// 计算时间
	elapsed := time.Since(start)

	// 发送计算结果
	communityOutput := graph.DisplayCommunities()
	fmt.Fprintf(conn, "%s", communityOutput)
	fmt.Fprintf(conn, "Temps d'exécution: %s\n", elapsed)
	fmt.Fprintln(conn, "FIN")
}

func main() {
	// 创建 TCP 监听
	listener, err := net.Listen("tcp", ":8882")
	if err != nil {
		fmt.Println("Erreur de démarrage du serveur :", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serveur en attente de connexions sur le port 8882...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur d'acceptation :", err)
			continue
		}
		go handleConnection(conn) // 并发处理多个客户端
	}
}
