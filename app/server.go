package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
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
	if err != nil {
		os.Exit(1)
	}
	path := request.URL.Path
	var response []byte

	switch {
	case path == "/":
		response = ([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case strings.HasPrefix(path, "/echo/"):
		encodings := request.Header.Get("Accept-Encoding")
		splitEncodings := strings.Split(encodings, ", ")
		message := path[6:]

		var b bytes.Buffer
		wr := gzip.NewWriter(&b)
		wr.Write([]byte(message))
		wr.Close()
		if slices.Contains(splitEncodings, "gzip") {
			response = ([]byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s",
				len(b.String()),
				b.String(),
			)))
		} else {
			response = ([]byte(fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
				len(message),
				message,
			)))

		}

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
