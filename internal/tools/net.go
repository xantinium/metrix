package tools

import (
	"net"
)

const (
	HeaderXRealIP         = "X-Real-IP"
	HeaderAccept          = "Accept"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentType     = "Content-Type"
	HeaderHashSHA256      = "HashSHA256"
)

// FindLocalIP находит локальный IP-адрес.
func FindLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, nil
}
