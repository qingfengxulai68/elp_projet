package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 连接到 localhost 的 8080 端口
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 创建一个扫描器来读取用户输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter numbers separated by commas (e.g., 2,3,9,7,1):")

	// 读取用户输入
	input, _ := reader.ReadString('\n')

	// 发送输入到服务端
	_, err = conn.Write([]byte(input))
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return
	}

	// 接收服务端的响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	// 打印服务端返回的结果
	fmt.Printf("Received from server: %s\n", string(buffer[:n]))
}
