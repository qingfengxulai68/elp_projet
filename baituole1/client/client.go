package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// 读取文件并发送数据
	file, err := os.Open("minigraph.txt")
	if err != nil {
		fmt.Println("Erreur d'ouverture du fichier :", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fmt.Fprintln(writer, line)
	}
	fmt.Fprintln(writer, "RUN") // 发送 RUN 命令
	writer.Flush()

	// 读取服务器返回的结果
	serverReader := bufio.NewScanner(conn)
	fmt.Println("Résultat de la détection de communauté :")
	for serverReader.Scan() {
		line := serverReader.Text()
		if line == "FIN" {
			break
		}
		fmt.Println(line)
	}
}
