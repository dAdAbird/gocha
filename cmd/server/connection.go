package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type cstate int

const (
	scInit cstate = 1 << iota
	scAuth
	scClosed
)

type connection struct {
	c           net.Conn
	state       cstate // auth, close, err, slow
	user        *user
	currChannel *Channel
	colse       chan struct{}
}

func (c *connection) Process(cmds *CommadsSet) {
	reader := bufio.NewReader(c.c)
	for {
		select {
		case <-c.colse:
			// fmt.Printf("Closing: %v", c)
			return
		default:
			msg, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println("stream read error:", err)
				_, err := CmdQuit(c, nil)
				log.Println("closing connection:", err)
			}

			// fmt.Printf("debug: %s\n", msg)

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

func (c *connection) SendString(msg string) (int, error) {
	return c.Send([]byte(msg))
}

func (c *connection) Send(msg []byte) (int, error) {
	// TODO Check errors and close connection
	return c.c.Write(msg)
}

type cPool struct {
	mx sync.Mutex
	c  map[*connection]struct{}
}

func (p *cPool) Add(c *connection) {
	p.mx.Lock()
	p.c[c] = struct{}{}
	p.mx.Unlock()
}

func (p *cPool) Len() int {
	p.mx.Lock()
	defer p.mx.Unlock()
	return len(p.c)
}

func (p *cPool) Delete(c *connection) int {
	p.mx.Lock()
	defer p.mx.Unlock()
	delete(p.c, c)
	return len(p.c)
}
