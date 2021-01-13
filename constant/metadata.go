package constant

import (
	"encoding/json"
	"net"
	"strconv"
)

// Socks addr type
const (
	ATypIPv4       = 1
	ATypDomainName = 3
	ATypIPv6       = 4

	TCP NetWork = iota
	UDP
)

type NetWork int

func (n NetWork) String() string {
	if n == TCP {
		return "tcp"
	}
	return "udp"
}

func (n NetWork) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

// Metadata is used to store connection address
type Metadata struct {
	NetWork  NetWork `json:"network"`
	DstIP    net.IP  `json:"destinationIP"`
	DstPort  string  `json:"destinationPort"`
	AddrType int     `json:"-"`
	Host     string  `json:"host"`
}

func (m *Metadata) RemoteAddress() string {
	return net.JoinHostPort(m.String(), m.DstPort)
}

func (m *Metadata) SourceAddress() string {
	return "nil"
}

func (m *Metadata) Resolved() bool {
	return m.DstIP != nil
}

func (m *Metadata) UDPAddr() *net.UDPAddr {
	if m.NetWork != UDP || m.DstIP == nil {
		return nil
	}
	port, _ := strconv.Atoi(m.DstPort)
	return &net.UDPAddr{
		IP:   m.DstIP,
		Port: port,
	}
}

func (m *Metadata) String() string {
	if m.Host != "" {
		return m.Host
	} else if m.DstIP != nil {
		return m.DstIP.String()
	} else {
		return "<nil>"
	}
}

func (m *Metadata) Valid() bool {
	return m.Host != "" || m.DstIP != nil
}
