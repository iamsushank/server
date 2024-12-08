package socket

import (
	"fmt"
	"syscall"
)

func FormatAddr(addr syscall.Sockaddr) string {
	if v4, ok := addr.(*syscall.SockaddrInet4); ok {
		return fmt.Sprintf("%d.%d.%d.%d:%d", v4.Addr[0], v4.Addr[1], v4.Addr[2], v4.Addr[3], v4.Port)
	}

	return "unknown address"
}
