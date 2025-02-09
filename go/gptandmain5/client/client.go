package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// 让用户选择文件编号
	var fileNumber, goroutineCount int
	fmt.Print("Veuillez choisir un fichier (1-3) : ")
	_, err := fmt.Scanln(&fileNumber)
	if err != nil || fileNumber < 1 || fileNumber > 3 {
		fmt.Println("Numéro invalide ! Veuillez entrer 1, 2 ou 3.")
		return
	}

	// 让用户输入 Goroutine 数量
	fmt.Print("Veuillez entrer le nombre de Goroutines : ")
	_, err = fmt.Scanln(&goroutineCount)
	if err != nil || goroutineCount < 1 {
		fmt.Println("Nombre invalide ! Veuillez entrer un nombre positif.")
		return
	}

	// 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8882")
	if err != nil {
		fmt.Println("Erreur de connexion au serveur :", err)
		return
	}
	defer conn.Close()

	// 发送编号和 Goroutine 数量给服务器
	writer := bufio.NewWriter(conn)
	fmt.Fprintf(writer, "%d %d\n", fileNumber, goroutineCount)
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
