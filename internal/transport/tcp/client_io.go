package tcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (c *Client) readFromServer() {
	buf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			fmt.Println("Disconnected")
			os.Exit(0)
		}

		ClearScreen()

		msg := string(buf[:n])
		fmt.Print(msg)

		if !strings.HasSuffix(msg, "\n") {
			fmt.Println()
		}

		fmt.Print("> ")
	}
}

func (c *Client) writeToServer() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		_, _ = c.conn.Write([]byte(text + "\n"))
	}
}
