package main

import (
	"baituole1/traitement"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// var wg sync.WaitGroup
var semaphore = make(chan struct{}, 2) // 限制最多 20个并发执行 Louvain

// var semaphore = make(chan struct{}, runtime.NumCPU()) // 动态设置 Goroutine 限制

// 处理客户端请求
func handleConnection(conn net.Conn, g *traitement.Graph) {

	defer conn.Close()

	// 限制Goroutine数量
	semaphore <- struct{}{}        // 获取 Goroutine 许可
	defer func() { <-semaphore }() // 释放 Goroutine 许可

	// // 创建一个 sync.Mutex 实例
	// var mu sync.Mutex

	scanner := bufio.NewScanner(conn)

	// 读取图数据
	for scanner.Scan() {
		line := scanner.Text()
		if line == "RUN" {
			break // 当收到 RUN 指令时，开始计算 Louvain
		}

		nodes := strings.Fields(line)
		if len(nodes) == 2 {
			u, _ := strconv.Atoi(nodes[0])
			v, _ := strconv.Atoi(nodes[1])
			g.AddEdge(u, v)
		}
	}

	// 添加 Mutex 保护，确保并发操作时的 map 访问是安全的
	var mu sync.Mutex

	// 并行运行 Louvain
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Louvain(10, &mu)
	}()
	wg.Wait()

	// 发送结果给客户端
	for node, community := range g.Communities {
		fmt.Fprintf(conn, "%d -> Community %d\n", node, community)
	}
	fmt.Fprintln(conn, "FIN") // 结束标志
}

func main() {
	// 创建 TCP 监听
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erreur de démarrage du serveur :", err)
		return
	}
	defer listener.Close()

	fmt.Println("Serveur en attente de connexions sur le port 8080...")

	graph := traitement.NewGraph()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur d'acceptation :", err)
			continue
		}
		go handleConnection(conn, graph) // 启动协程处理连接
	}
}
