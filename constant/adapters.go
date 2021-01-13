package constant

import (
	"context"
	"fmt"
	"net"
)

// Adapter Type
const (
	Shadowsocks AdapterType = iota
	ShadowsocksR
	Snell
	Socks5
	Http
	Vmess
	Trojan
)

type ServerAdapter interface {
	net.Conn
	Metadata() *Metadata
}

type Connection interface {
	Chains() Chain
	AppendToChains(adapter ProxyAdapter)
}

type Chain []string

func (c Chain) String() string {
	switch len(c) {
	case 0:
		return ""
	case 1:
		return c[0]
	default:
		return fmt.Sprintf("%s[%s]", c[len(c)-1], c[0])
	}
}

type Conn interface {
	net.Conn
	Connection
}

type PacketConn interface {
	net.PacketConn
	Connection
	// Deprecate WriteWithMetadata because of remote resolve DNS cause TURN failed
	// WriteWithMetadata(p []byte, metadata *Metadata) (n int, err error)
}

type ProxyAdapter interface {
	Name() string
	Type() AdapterType
	StreamConn(c net.Conn, metadata *Metadata) (net.Conn, error)
	DialContext(ctx context.Context, metadata *Metadata) (Conn, error)
	DialUDP(metadata *Metadata) (PacketConn, error)
	SupportUDP() bool
	Addr() string
}

type Proxy interface {
	ProxyAdapter
	Dial(metadata *Metadata) (Conn, error)
	URLTest(ctx context.Context, URL string) (uint16, error)
}

// AdapterType is enum of adapter type
type AdapterType int

func (at AdapterType) String() string {
	switch at {
	case Shadowsocks:
		return "Shadowsocks"
	case ShadowsocksR:
		return "ShadowsocksR"
	case Snell:
		return "Snell"
	case Socks5:
		return "Socks5"
	case Http:
		return "Http"
	case Vmess:
		return "Vmess"
	case Trojan:
		return "Trojan"

	default:
		return "Unknown"
	}
}

// UDPPacket contains the data of UDP packet, and offers control/info of UDP packet's source
type UDPPacket interface {
	// Data get the payload of UDP Packet
	Data() []byte

	// WriteBack writes the payload with source IP/Port equals addr
	// - variable source IP/Port is important to STUN
	// - if addr is not provided, WriteBack will write out UDP packet with SourceIP/Port equals to original Target,
	//   this is important when using Fake-IP.
	WriteBack(b []byte, addr net.Addr) (n int, err error)

	// Drop call after packet is used, could recycle buffer in this function.
	Drop()

	// LocalAddr returns the source IP/Port of packet
	LocalAddr() net.Addr
}
