package tcp

import "net"

type ClientConn struct {
	Name      string
	SessionID string
	Conn      net.Conn
}
