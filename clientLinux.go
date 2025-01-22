package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 连接到 localhost 的 8080 端口
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 提示用户输入文件路径
	fmt.Print("Enter the path to the .txt file: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	// 发送文件路径到服务端
	_, err = conn.Write([]byte(filePath))
	if err != nil {
		fmt.Println("Error sending file path to server:", err)
		return
	}

	// 接收服务端的确认消息
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from server:", err)
		return
	}

	// 打印服务端返回的结果
	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
