package tcp

import "net"

type ClientConn struct {
	Name string
	Conn net.Conn
}
