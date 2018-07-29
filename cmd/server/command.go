package main

import (
	"bytes"
	"errors"
)

type cmdAction func(c *connection, b []byte) (resp string, err error)

func cmdWrapAuth(f cmdAction) cmdAction {
	return func(c *connection, msg []byte) (resp string, err error) {
		if c.state&scAuth == 0 {
			return "", ErrUnauthUser
		}
		return f(c, msg)
	}
}

type CommadsSet map[string]cmdAction

func (cs CommadsSet) Register(name string, action cmdAction) {
	cs[name] = action
}

var ErrNoCommandSet = errors.New("no command set")

func (cs CommadsSet) Parse(msg []byte) (cmd cmdAction, args []byte, err error) {
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
	return nil, nil, ErrNoCommandSet
}

func CmdSend(c *connection, msg []byte) (resp string, err error) {
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

func CmdQuit(c *connection, _ []byte) (resp string, err error) {
	if c.state&scAuth != 0 {
		left := c.user.conns.Delete(c)
		if left == 0 {
			c.currChannel.UnRegister(c.user)
		}

		c.SendString("BYE!\n")
	}

	c.state |= scClosed
	c.colse <- struct{}{}
	return resp, c.c.Close()
}

var ErrUnauthUser = errors.New("user unknown")

func CmdAuth(c *connection, name []byte) (resp string, err error) {
	if len(name) == 0 {
		return "", ErrUserNoName
	}

	if c.state&scAuth == 0 {
		c.user = connectedUsers.Add(string(name), c)

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
