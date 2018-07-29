package main

import (
	"bytes"
	"errors"
)

type cmdAction func(c *connection, b []byte) (resp string, err error)
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

func RespMsg(c *connection, msg []byte) (resp string, err error) {
	// c.cliWriter.WriteByte('\n')
	c.cliWriter.Write(msg)
	c.cliWriter.WriteByte('\n')

	// c.pipe <- struct{}{}
	// return
	return "", c.cliWriter.Flush()
}

func RespError(c *connection, msg []byte) (resp string, err error) {
	c.cliWriter.WriteString("ERROR: ")
	c.cliWriter.Write(msg)

	// c.pipe <- struct{}{}
	// return
	return "", c.cliWriter.Flush()
}

func RespOK(c *connection, msg []byte) (resp string, err error) {
	if c.lastCmd == "AUTH" {
		c.cliWriter.WriteString("*** Hi " + c.userName)
		if i := bytes.IndexByte(msg, ':'); i > 0 {
			c.currentChannel = string(msg[:i])
			c.cliWriter.WriteString(", here is")
			c.cliWriter.Write(msg[i+1:])
			c.cliWriter.WriteByte('.')
			c.cliWriter.WriteByte('\n')
		} else {
			c.currentChannel = string(msg)
			c.cliWriter.WriteString(", there is nobody here yet.")
			c.cliWriter.WriteByte('\n')
		}

		c.cliWriter.Flush()
	}

	// c.pipe <- struct{}{}

	return
}
