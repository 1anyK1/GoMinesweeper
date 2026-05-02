package tcp

import (
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) Run() {
	go c.readFromServer()
	c.writeToServer()
}
