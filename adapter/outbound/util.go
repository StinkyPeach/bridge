package outbound

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/StinkyPeach/bridge/component/resolver"
	C "github.com/StinkyPeach/bridge/constant"
)

const (
	tcpTimeout = 5 * time.Second
)

var (
	globalClientSessionCache tls.ClientSessionCache
	once                     sync.Once
)

func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	addr = C.Metadata{
		AddrType: C.ATypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return
}

func tcpKeepAlive(c net.Conn) {
	if tcp, ok := c.(*net.TCPConn); ok {
		tcp.SetKeepAlive(true)
		tcp.SetKeepAlivePeriod(30 * time.Second)
	}
}

func getClientSessionCache() tls.ClientSessionCache {
	once.Do(func() {
		globalClientSessionCache = tls.NewLRUClientSessionCache(128)
	})
	return globalClientSessionCache
}

func serializesSocksAddr(metadata *C.Metadata) []byte {
	var buf [][]byte
	aType := uint8(metadata.AddrType)
	p, _ := strconv.Atoi(metadata.DstPort)
	port := []byte{uint8(p >> 8), uint8(p & 0xff)}
	switch metadata.AddrType {
	case C.ATypDomainName:
		len := uint8(len(metadata.Host))
		host := []byte(metadata.Host)
		buf = [][]byte{{aType, len}, host, port}
	case C.ATypIPv4:
		host := metadata.DstIP.To4()
		buf = [][]byte{{aType}, host, port}
	case C.ATypIPv6:
		host := metadata.DstIP.To16()
		buf = [][]byte{{aType}, host, port}
	}
	return bytes.Join(buf, nil)
}

func resolveUDPAddr(network, address string) (*net.UDPAddr, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	ip, err := resolver.ResolveIP(host)
	if err != nil {
		return nil, err
	}
	return net.ResolveUDPAddr(network, net.JoinHostPort(ip.String(), port))
}
