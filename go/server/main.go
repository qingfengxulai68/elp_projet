package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// AddEdge 向图中添加一条边
func (g *Graph) AddEdge(node1, node2 string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.AdjList[node1] = append(g.AdjList[node1], node2)
	g.AdjList[node2] = append(g.AdjList[node2], node1) // 无向图
}

// Print 打印图的邻接表
func (g *Graph) Print() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for node, neighbors := range g.AdjList {
		fmt.Printf("%s: %v\n", node, neighbors)
	}
}

// LoadGraphFromFile 从指定路径加载图
func LoadGraphFromFile(filePath string) (*Graph, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := &Graph{AdjList: make(map[string][]string)}
	reader := bufio.NewReader(file)

	// 使用 bufio.Reader 逐行读取文件，处理长行
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			break
		}
		line = strings.TrimSpace(line)
		if line != "" {
			// 修改这里：使用制表符 `\t` 分割行
			nodes := strings.Split(line, "\t")
			if len(nodes) == 2 {
				graph.AddEdge(nodes[0], nodes[1]) // 添加边
			}
		}
	}

	return graph, nil
}

// handleConnection 处理每个客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected.")

	// 接收客户端传递的文件路径
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	// 解析客户端发送的文件路径
	filePath := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("Received file path from client: %s\n", filePath)

	// 如果客户端传递的是 Windows 风格路径，需要转换为 Linux 路径
	if runtime.GOOS == "linux" && strings.HasPrefix(filePath, "C:") {
		// 转换 Windows 路径为 Linux 路径，确保不多出冒号
		filePath = "/mnt/c" + filePath[2:]
		// 替换反斜杠 `\` 为正斜杠 `/`
		filePath = strings.ReplaceAll(filePath, "\\", "/")
	}

	// 加载图结构
	graph, err := LoadGraphFromFile(filePath)
	if err != nil {
		fmt.Fprintf(conn, "Error loading graph from file: %v\n", err)
		fmt.Println("Error loading graph:", err)
		return
	}

	// 打印图结构到服务器控制台
	fmt.Println("Graph loaded successfully. Adjacency list:")
	graph.PrintGraph()

	// 向客户端确认加载成功
	_, err = conn.Write([]byte("Graph loaded successfully. Check server logs for adjacency list.\n"))
	if err != nil {
		fmt.Println("Error sending confirmation to client:", err)
		return
	}

	fmt.Println("Graph loading confirmation sent to client.")

}

func algo() {
	// Créer un nouveau graphe
	graph := NewGraph()

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

	// Temps d'exécution
	start := time.Now()

	// Lancer l'algorithme de Louvain pour détecter les communautés
	graph.Louvain(10)
	fmt.Println("Communautés trouvées :")
	// Afficher les communautés
	graph.DisplayCommunities()

	// Afficher le temps d'exécution
	elapsed := time.Since(start)
	fmt.Printf("Temps d'exécution: %s\n", elapsed)
}

// main 启动 TCP 服务器
func main() {
	// 监听本地的 8080 端口
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	// 等待客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 启动一个 Goroutine 处理每个客户端连接
		go handleConnection(conn)
	}
}
