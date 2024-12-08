package main

import (
	"server/internal/http"
	"server/internal/server"
)

func main() {
	port := 8080

	// Register routes
	server.RegisterRoute("/static", func(request *http.Request) string {
		return http.GenerateResponse(200, `{"message": "Hello, this is static data!"}`)
	})

	server.RegisterRoute("/hello", func(request *http.Request) string {
		return http.GenerateResponse(200, `{"message": "Welcome to my custom HTTP server!"}`)
	})

	// Start the server
	server.StartServer(port)
}
