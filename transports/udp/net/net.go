package net

import (
	"net"
	"os"

	honeybadger "github.com/honeybadger-io/honeybadger-go"
)

var ResolveUDPAddr = net.ResolveUDPAddr

func DialUDP(n string, laddr, raddr *net.UDPAddr) (*UDPConn, error) {
	conn, err := net.DialUDP(n, laddr, raddr)
	return &UDPConn{UDPConn: conn}, err
}

// UDPConn wraps net.UDPConn with error reports. The pimary reason we have this
// in a `net` package with the name UDPConn is because the syslog adapter uses
// reflecting to look at the connection type. We want to make sure that it looks
// like an actual net.UDPConn.
type UDPConn struct {
	*net.UDPConn
	HB *honeybadger.Client
}

func (conn *UDPConn) Write(b []byte) (int, error) {
	n, err := conn.UDPConn.Write(b)
	if err != nil {
		if os.Getenv("HONEYBADGER_API_KEY") != "" {
			conn.hb().Notify(err, honeybadger.Context{
				"data": string(b),
			})
		}
	}
	return n, err
}

func (conn *UDPConn) hb() *honeybadger.Client {
	if conn.HB == nil {
		return honeybadger.DefaultClient
	}

	return conn.HB
}
