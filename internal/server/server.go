package server

import (
	"fmt"
	"server/internal/http"
	"server/internal/socket"
	"syscall"
)

var routes = map[string]func(*http.Request) string{}

func RegisterRoute(path string, handler func(*http.Request) string) {
	routes[path] = handler
}

func StartServer(port int) {
	fd, err := socket.CreateSocket()
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)

	if err := socket.BindSocket(fd, port); err != nil {
		panic(err)
	}

	if err := socket.ListenSocket(fd); err != nil {
		panic(err)
	}

	fmt.Printf("Server is listening on port %d\n", port)

	for {
		connFd, addr, err := socket.AcceptConnection(fd)
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(connFd, addr)
	}
}

func handleConnection(connFd int, addr syscall.Sockaddr) {
	defer syscall.Close(connFd)

	buffer := make([]byte, 4096)
	n, err := syscall.Read(connFd, buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	requestData := string(buffer[:n])
	request, err := http.ParseRequest(requestData)
	if err != nil {
		fmt.Printf("Invalid request: %v\n", err)
		response := http.GenerateResponse(400, "Bad Request")
		syscall.Write(connFd, []byte(response))
		return
	}

	response := routeRequest(request)
	syscall.Write(connFd, []byte(response))
}

func routeRequest(request *http.Request) string {
	if handler, exists := routes[request.URL]; exists {
		return handler(request)
	}
	return http.GenerateResponse(404, "Not Found")
}
