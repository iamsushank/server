package socket

import (
	"fmt"
	"syscall"
)

func CreateSocket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return -1, fmt.Errorf("failed to create socket: %w", err)
	}

	return fd, nil
}

func BindSocket(fd, port int) error {
	addr := &syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], []byte{127, 0, 0, 1})

	if err := syscall.Bind(fd, addr); err != nil {
		return fmt.Errorf("failed to bind socket: %w", err)
	}

	return nil
}

func ListenSocket(fd int) error {
	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	return nil
}

func AcceptConnection(fd int) (int, syscall.Sockaddr, error) {
	connFd, addr, err := syscall.Accept(fd)
	if err != nil {
		return -1, nil, fmt.Errorf("failed to accept connection: %w", err)
	}

	return connFd, addr, nil
}
