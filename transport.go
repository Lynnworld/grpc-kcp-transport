package kcp

import (
	"context"
	"fmt"
	xkcp "github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/tcpraw"
	"google.golang.org/grpc"
	"net"
)

// Dialer for grpc client dial
func Dialer(c *Config) grpc.DialOption {
	c.fixture()

	fun := func(ctx context.Context, addr string) (kcp net.Conn, err error) {
		block, err := getCryptBlock(c.Crypt, c.Key)
		if err != nil {
			return nil, err
		}

		var kcpconn *xkcp.UDPSession
		if c.TCP {
			conn, err := tcpraw.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			kcpconn, err = xkcp.NewConn(addr, block, c.DataShard, c.ParityShard, conn)
		} else {
			kcpconn, err = xkcp.DialWithOptions(addr, block, c.DataShard, c.ParityShard)
		}

		if err != nil {
			return nil, err
		}
		kcpconn.SetStreamMode(true)
		kcpconn.SetWriteDelay(false)

		kcpconn.SetNoDelay(c.NoDelay, c.Interval, c.Resend, c.NoCongestion)
		kcpconn.SetWindowSize(c.SndWnd, c.RcvWnd)
		kcpconn.SetMtu(c.MTU)
		kcpconn.SetACKNoDelay(c.AckNodelay)

		_ = kcpconn.SetDSCP(c.DSCP)
		_ = kcpconn.SetReadBuffer(c.SockBuf)
		_ = kcpconn.SetWriteBuffer(c.SockBuf)

		return kcpconn, err
	}
	return grpc.WithContextDialer(fun)
}

type listener struct {
	*xkcp.Listener
	c *Config
}

func (l *listener) Accept() (net.Conn, error) {
	conn, err := l.Listener.AcceptKCP()
	if err != nil {
		return nil, err
	}
	conn.SetStreamMode(true)
	conn.SetWriteDelay(false)
	conn.SetNoDelay(l.c.NoDelay, l.c.Interval, l.c.Resend, l.c.NoCongestion)
	conn.SetMtu(l.c.MTU)
	conn.SetWindowSize(l.c.SndWnd, l.c.RcvWnd)
	conn.SetACKNoDelay(l.c.AckNodelay)
	return conn, nil
}

// ServeGrpc for grpc server
func ServeGrpc(address string, server *grpc.Server, c *Config) error {
	// tcp mode
	if c.TCP {
		lis, err := Listener(address, c)
		if err != nil {
			return err
		}
		return server.Serve(lis)
	}
	// udp mode
	laddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	if !laddr.IP.IsUnspecified() {
		listen, err := Listener(address, c)
		if err != nil {
			return err
		}
		return server.Serve(listen)
	}
	ipList := getLocalAddress()
	localIp := ipList[0]
	ipList = ipList[1:]
	for _, ip := range ipList {
		addr := fmt.Sprintf("%s:%d", ip, laddr.Port)
		listen, err := Listener(addr, c)
		if err != nil {
			return err
		}
		go func() {
			_ = server.Serve(listen)
		}()
	}
	addr := fmt.Sprintf("%s:%d", localIp, laddr.Port)
	listen, err := Listener(addr, c)
	if err != nil {
		return err
	}
	return server.Serve(listen)
}

func getLocalAddress() []string {
	ips := make([]string, 0)
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

// Listener Listen for grpc server listen, addr is like "127.0.0.1:8080"
func Listener(address string, c *Config) (l net.Listener, err error) {
	c.fixture()
	block, err := getCryptBlock(c.Crypt, c.Key)
	if err != nil {
		return nil, err
	}

	var ln *xkcp.Listener
	if c.TCP {
		var conn *tcpraw.TCPConn
		conn, err = tcpraw.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
		ln, err = xkcp.ServeConn(block, c.DataShard, c.ParityShard, conn)
	} else {
		ln, err = xkcp.ListenWithOptions(address, block, c.DataShard, c.ParityShard)
	}
	if err != nil {
		return nil, err
	}
	_ = ln.SetDSCP(c.DSCP)
	_ = ln.SetReadBuffer(c.SockBuf)
	_ = ln.SetWriteBuffer(c.SockBuf)

	l = &listener{
		Listener: ln,
		c:        c,
	}
	return
}
