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

	// 指定需要读取的文件路径
	filePath := "path/to/your/file.txt" // 替换为你的实际文件路径

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个扫描器读取文件内容
	scanner := bufio.NewScanner(file)
	var input string

	for scanner.Scan() {
		// 拼接文件的每一行内容，并用换行符分隔
		input += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 输出读取到的文件内容（可用于调试）
	fmt.Println("Read from file:")
	fmt.Println(input)

	// 发送文件内容到服务端
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
