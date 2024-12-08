package http

import "fmt"

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func GenerateResponse(statusCode int, body string) string {
	statusText := map[int]string{
		200: "OK",
		404: "Not Found",
		500: "Internal Server Error",
	}

	headers := fmt.Sprintf("Content-Length: %d\r\n", len(body))
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", statusCode, statusText[statusCode], headers, body)
}
