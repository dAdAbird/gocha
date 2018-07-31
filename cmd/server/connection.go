package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

type cstate int

const (
	scInit cstate = 1 << iota
	scAuth
	scClosed
)

// Client represents the client connection
type Client struct {
	conn        net.Conn
	state       cstate // auth, closed etc.
	user        *User
	currChannel *Channel
	colse       chan error
}

// Handle process the client connection
func (c *Client) Handle(cmds *CommadsSet) {
	reader := bufio.NewReader(c.conn)
	for {
		select {
		case cerr := <-c.colse:
			if cerr != nil {
				log.Println("Closing connection because of:", cerr)
			}
			_, err := cmdQuit(c, nil)
			log.Println("Connection closed:", err)
			return
		default:
			msg, err := reader.ReadBytes('\n')
			if err != nil {
				c.colse <- fmt.Errorf("stream read error: %v", err)
			}

			c.conn.SetDeadline(time.Now().Add(15 * time.Minute))

			runCmd, args, err := cmds.Parse(msg)
			if err != nil {
				c.SendString(fmt.Sprintf("ERR %v\n", err))
				continue
			}

			resp, err := runCmd(c, args)
			if err != nil {
				c.SendString(fmt.Sprintf("ERR %v\n", err))
				continue
			}

			c.SendString("OK " + resp + "\n")
		}
	}
}

// SendString msg to client
func (c *Client) SendString(msg string) (int, error) {
	return c.Send([]byte(msg))
}

// Send msg to client
func (c *Client) Send(msg []byte) (int, error) {
	return c.conn.Write(msg)
}
