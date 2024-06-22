package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		os.Exit(1)
	}
	path := request.URL.Path
	var response []byte

	switch {
	case path == "/":
		response = ([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.HasPrefix(path, "/echo/"):
		message := path[6:]
		response = ([]byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(message),
			message,
		)))

	case path == "/user-agent":
		response = ([]byte(fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
			len(request.UserAgent()),
			request.UserAgent(),
		)))
	default:
		response = ([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	conn.Write(response)
}
