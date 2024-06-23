# HTTP Server in Go

This is a simple HTTP server implemented in Go. It listens for connections on a specified port and provides various HTTP endpoints to demonstrate basic HTTP server functionalities.

## Features

- **Root Endpoint (`/`)**: Returns a simple HTTP 200 OK response.
- **Echo Endpoint (`/echo/{message}`)**: Returns the message in the URL path. Supports gzip compression if requested by the client.
- **User-Agent Endpoint (`/user-agent`)**: Returns the User-Agent string of the client.
- **File Endpoints**: 
  - **GET /files/{filename}**: Serves files from the server's filesystem.
  - **POST /files/{filename}**: Allows uploading files to the server's filesystem.

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/nemopss/http-server-go.git
   ```

2. Change into the project directory:

   ```sh
   cd http-server-go
   ```

3. Build the server:

   ```sh
   go build -o http-server app/server.go
   ```

## Usage

To start the server, run the following command:

```sh
./http-server /path/to/files
```

This will start the server and bind it to port `4221`.

### Endpoints

- **Root Endpoint:**
  ```sh
  curl http://localhost:4221/
  ```

- **Echo Endpoint:**
  ```sh
  curl http://localhost:4221/echo/your_message
  ```

  To request gzip encoding:
  ```sh
  curl -H "Accept-Encoding: gzip" http://localhost:4221/echo/your_message
  ```

- **User-Agent Endpoint:**
  ```sh
  curl http://localhost:4221/user-agent
  ```

- **File Endpoints:**
  - **GET /files/{filename}**:
    ```sh
    curl http://localhost:4221/files/your_file.txt
    ```

  - **POST /files/{filename}**:
    ```sh
    curl -X POST --data-binary @local_file.txt http://localhost:4221/files/your_file.txt
    ```
