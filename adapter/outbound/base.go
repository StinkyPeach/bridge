package outbound

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	C "github.com/StinkyPeach/bridge/constant"
	"net"
	"net/http"
	"time"
)

type Base struct {
	name string
	addr string
	tp   C.AdapterType
	udp  bool
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Type() C.AdapterType {
	return b.tp
}

func (b *Base) StreamConn(c net.Conn, metadata *C.Metadata) (net.Conn, error) {
	return c, errors.New("no support")
}

func (b *Base) DialUDP(metadata *C.Metadata) (C.PacketConn, error) {
	return nil, errors.New("no support")
}

func (b *Base) SupportUDP() bool {
	return b.udp
}

func (b *Base) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"type": b.Type().String(),
	})
}

func (b *Base) Addr() string {
	return b.addr
}

func NewBase(name string, addr string, tp C.AdapterType, udp bool) *Base {
	return &Base{name, addr, tp, udp}
}

type conn struct {
	net.Conn
	chain C.Chain
}

func (c *conn) Chains() C.Chain {
	return c.chain
}

func (c *conn) AppendToChains(a C.ProxyAdapter) {
	c.chain = append(c.chain, a.Name())
}

func NewConn(c net.Conn, a C.ProxyAdapter) C.Conn {
	return &conn{c, []string{a.Name()}}
}

type packetConn struct {
	net.PacketConn
	chain C.Chain
}

func (c *packetConn) Chains() C.Chain {
	return c.chain
}

func (c *packetConn) AppendToChains(a C.ProxyAdapter) {
	c.chain = append(c.chain, a.Name())
}

func newPacketConn(pc net.PacketConn, a C.ProxyAdapter) C.PacketConn {
	return &packetConn{pc, []string{a.Name()}}
}

type Proxy struct {
	C.ProxyAdapter
}

func (p *Proxy) Dial(metadata *C.Metadata) (C.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tcpTimeout)
	defer cancel()
	return p.DialContext(ctx, metadata)
}

func (p *Proxy) DialContext(ctx context.Context, metadata *C.Metadata) (C.Conn, error) {
	conn, err := p.ProxyAdapter.DialContext(ctx, metadata)

	return conn, err
}

// URLTest get the delay for the specified URL, t ms
func (p *Proxy) URLTest(ctx context.Context, URL string) (t uint16, err error) {
	metadata, err := urlToMetadata(URL)
	if err != nil {
		return 0, err
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := p.Dial(&metadata)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	client := http.Client{
		Transport: &http.Transport{
			DialContext: dialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: tcpTimeout,
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	sTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	t = uint16(time.Since(sTime).Milliseconds())

	return t, nil
}

func NewProxy(adapter C.ProxyAdapter) *Proxy {
	return &Proxy{adapter}
}
