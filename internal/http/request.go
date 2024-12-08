package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
	Body    string
}

func ParseRequest(data string) (*Request, error) {
	reader := bufio.NewReader(strings.NewReader(data))

	request := &Request{
		Headers: make(map[string]string),
	}

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read request line: %w", err)
	}

	requestLine = strings.TrimSpace(requestLine)

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	request.Method, request.URL, request.Version = parts[0], parts[1], parts[2]

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read headers: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" { // Empty line indicates end of headers
			break
		}

		parts = strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}

		request.Headers[parts[0]] = parts[1]
	}

	if transferEncoding, ok := request.Headers["Transfer-Encoding"]; ok && transferEncoding == "chunked" {
		body, err := parseChunkedBody(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chunked body: %w", err)
		}
		request.Body = body
	} else {
		body, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		request.Body = string(body)
	}

	return request, nil
}

func parseChunkedBody(reader *bufio.Reader) (string, error) {
	var body bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read chunk size: %w", err)
		}
		line = strings.TrimSpace(line)
		chunkSize, err := strconv.ParseInt(line, 16, 64)
		if err != nil {
			return "", fmt.Errorf("invalid chunk size: %w", err)
		}
		if chunkSize == 0 {
			break
		}

		chunk := make([]byte, chunkSize)
		_, err = io.ReadFull(reader, chunk)
		if err != nil {
			return "", fmt.Errorf("failed to read chunk data: %w", err)
		}
		body.Write(chunk)

		// Read and discard the trailing CRLF
		_, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read CRLF after chunk: %w", err)
		}
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read trailer: %w", err)
		}
		if strings.TrimSpace(line) == "" {
			break
		}
	}

	return body.String(), nil
}
