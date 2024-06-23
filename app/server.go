package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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
		go func() { handleConnection(conn) }()

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	request, err := http.ReadRequest(bufio.NewReader(conn))
	fmt.Print(request.Method)
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

	case request.Method == "GET" && strings.HasPrefix(path, "/files/"):
		fileName := strings.TrimPrefix(path, "/files/")
		if bytes, err := os.ReadFile(os.Args[2] + fileName); err == nil {
			response = ([]byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
				len(bytes),
				string(bytes),
			)))
		} else {
			response = ([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	case request.Method == "POST" && strings.HasPrefix(path, "/files/"):
		fileName := strings.TrimPrefix(path, "/files/")
		length, _ := strconv.Atoi(request.Header.Get("Content-Length"))
		content := make([]byte, length)
		request.Body.Read(content)
		if err := os.WriteFile(os.Args[2]+fileName, content, 0644); err == nil {
			response = ([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		} else {

			response = ([]byte("HTTP/1.1 201 Created\r\n\r\n"))

		}
	default:
		response = ([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	conn.Write(response)
}
