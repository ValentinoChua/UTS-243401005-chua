package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Penggunaan: go run ./bukti-tcp <alamat>")
		return
	}

	alamat := os.Args[1]

	conn, err := net.Dial("tcp", alamat)
	if err != nil {
		fmt.Println("Gagal terhubung:", err)
		return
	}
	defer conn.Close()

	request := "GET /healthz HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Connection: close\r\n" +
		"\r\n"

	fmt.Fprint(conn, request)

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		fmt.Print(line)
	}
}
