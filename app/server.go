package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting the server...")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	req := make([]byte, 1024)
	conn.Read(req)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	path := strings.Split(string(req), " ")
	switch {
	case strings.HasPrefix(string(req), "GET / HTTP/1.1"):
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.HasPrefix(path[1], "/echo"):
		message := strings.Split(path[1], "/")[2]
		conn.Write([]byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(message),
			message)),
		)

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
