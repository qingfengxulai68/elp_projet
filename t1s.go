package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// findMax 计算最大值的函数
func findMax(numbers []int, wg *sync.WaitGroup, result chan<- int) {
	defer wg.Done() // 确保 Goroutine 完成时调用 Done()
	fmt.Println("Goroutine started, calculating max value...")

	max := numbers[0]
	for _, num := range numbers[1:] {
		if num > max {
			max = num
		}
	}

	// 将最大值发送到通道
	fmt.Println("Sending max value to channel...")
	result <- max // 这里发送数据到通道
	fmt.Println("Goroutine finished, result sent to channel.")
}

// handleConnection 处理每个客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected.")

	// 接收客户端发送的消息
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}

	// 将接收到的数据转换成字符串并去除多余的空格
	data := string(buffer[:n])
	fmt.Printf("Received data: %s\n", data)

	// 清理输入数据并转换成整数数组
	data = strings.TrimSpace(data)
	strNumbers := strings.Split(data, ",")
	numbers := make([]int, 0)

	// 处理每个字符串，转换为整数
	for _, str := range strNumbers {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Error converting to integer:", err)
			return
		}
		numbers = append(numbers, num)
	}

	// 调试输出：打印数字数组
	fmt.Printf("Converted numbers: %v\n", numbers)

	// 创建 WaitGroup 和通道来并发计算
	var wg sync.WaitGroup
	result := make(chan int, 1) // 设置缓冲区为1，确保至少有一个空间来接收数据

	// 调试输出：在添加 Goroutine 到 WaitGroup 时
	fmt.Println("Adding Goroutine to WaitGroup...")
	wg.Add(1)

	// 使用 Goroutine 来并行计算数组的最大值
	go findMax(numbers, &wg, result)

	// 打印 WaitGroup 大小（间接方式）
	fmt.Println("Waiting for Goroutine to complete...")

	// 等待所有 Goroutine 完成
	wg.Wait() // 等待 Goroutine 执行完毕
	fmt.Println("Goroutine completed!")

	// 从通道获取最大值并发送给客户端
	fmt.Println("Sending result to client...")
	close(result) // 关闭通道

	// 从通道读取数据并返回结果给客户端
	max := <-result // 等待从通道接收最大值
	fmt.Fprintf(conn, "The max value is: %d\n", max)

	// 发送结果后打印日志
	fmt.Println("Result sent to client.")
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

	fmt.Println("Server listening on port 8080 ")

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
