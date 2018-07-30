package main

import (
	"bytes"
	"errors"
)

// RespHandler represents the handler for command
type RespHandler func(s *Session, b []byte) error

// ResponsesSet contaign registred commands
type ResponsesSet map[string]RespHandler

// Register adds a new command into set
func (cs ResponsesSet) Register(name string, action RespHandler) {
	cs[name] = action
}

// ErrNoCommand error for undefined commans
var ErrNoCommand = errors.New("no command found")

// Parse parses the command from msg and return command Handler and Args
func (cs ResponsesSet) Parse(msg []byte) (cmd RespHandler, args []byte, err error) {
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
func respMsg(c *Session, msg []byte) error {
	c.cliWriter.WriteByte('\n')
	c.cliWriter.Write(msg)
	c.cliWriter.WriteByte('\n')

	c.cliWriter.WriteString(c.bar())

	return c.cliWriter.Flush()
}

func respError(c *Session, msg []byte) error {
	c.cliWriter.WriteByte('\n')
	c.cliWriter.WriteString("ERROR: ")
	c.cliWriter.Write(msg)

	c.cliWriter.WriteString(c.bar())

	return c.cliWriter.Flush()
}

func respOK(c *Session, msg []byte) error {
	if c.lastCmd == "AUTH" {
		c.cliWriter.WriteString("\n*** Hi ")
		c.cliWriter.WriteString(c.login)
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
		c.cliWriter.WriteString(c.bar())
		return c.cliWriter.Flush()
	}
	return nil
}
