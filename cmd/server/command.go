package main

import (
	"bytes"
	"errors"
	"fmt"
)

// CmdHandler represents the handler for command
type CmdHandler func(c *Client, b []byte) (resp string, err error)

func cmdWrapAuth(f CmdHandler) CmdHandler {
	return func(c *Client, msg []byte) (resp string, err error) {
		if c.state&scAuth == 0 {
			return "", ErrUnauthUser
		}
		return f(c, msg)
	}
}

// CommadsSet contaign registred commands
type CommadsSet map[string]CmdHandler

// Register adds a new command into set
func (cs CommadsSet) Register(name string, action CmdHandler) {
	cs[name] = action
}

// ErrNoCommand error for undefined commans
var ErrNoCommand = errors.New("no command found")

// Parse parses the command from msg and return command Handler and Args
func (cs CommadsSet) Parse(msg []byte) (cmd CmdHandler, args []byte, err error) {
	msg = bytes.TrimSpace(msg)
	cmdName := msg
	pos := bytes.IndexByte(msg, ' ')
	if pos > 0 {
		cmdName = msg[:pos]
	}

	if run, ok := cs[string(cmdName)]; ok {
		if len(msg) > len(cmdName) {
			// cut the space
			args = msg[pos+1:]
		}
		return run, args, nil
	}
	return nil, nil, ErrNoCommand
}

// Commands' handler fnctions.
// It sould be associated with text command and
// registred in CommadsSet
func cmdSend(c *Client, msg []byte) (resp string, err error) {
	send := bytes.NewBuffer(nil)
	send.Grow(4 + len(c.user.name) + len(msg) + 3) // :\t\n
	send.WriteString("MSG ")
	send.WriteString(c.user.name)
	send.WriteByte(':')
	send.WriteByte('\t')
	send.Write(msg)
	send.WriteByte('\n')

	c.currChannel.Send(send.Bytes(), c)
	return "", nil
}

func cmdQuit(c *Client, _ []byte) (resp string, err error) {
	if c.state&scAuth != 0 {
		left := c.user.conns.Delete(c)
		if left == 0 {
			c.currChannel.UnRegister(c.user)
		}

		c.SendString("BYE!\n")
	}

	c.state |= scClosed
	c.colse <- fmt.Errorf("recieved QUIT from client")
	return resp, c.conn.Close()
}

// ErrUnauthUser - user unauthorised
var ErrUnauthUser = errors.New("user unknown")

func cmdAuth(c *Client, name []byte) (resp string, err error) {
	if len(name) == 0 {
		return "", errors.New("no user name")
	}

	if c.state&scAuth == 0 {
		c.user = NewUser(string(name))
		// c.user.conns.Add(c)
		c.user.AddClient(c)

		resp = c.currChannel.name + ": "
		for _, u := range c.currChannel.Users() {
			resp += u.name + ", "
		}

		c.currChannel.Register(c.user)
		c.state |= scAuth

		if len(resp) > 2 {
			resp = resp[:len(resp)-2]
		}
		return resp, nil
	}
	return "", errors.New("that's not your name")
}
