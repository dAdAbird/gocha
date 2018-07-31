package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

// Session represents client chat Session
type Session struct {
	conn  net.Conn
	sHost string

	// chat data (login, current channel)
	login          string
	currentChannel string

	lastCmd string

	close chan struct{}

	netReader *bufio.Reader
	netWriter *bufio.Writer
	cliReader *bufio.Reader
	cliWriter *bufio.Writer
}

// Process contaigns main logic of chat client and process both
// cli input/output and network communication with server
func (s *Session) Process(cmds *ResponsesSet) {
	// starting go-routine to accept data from server
	go func() {
		for {
			select {
			case <-s.close:
				s.Sutdown()
				return
			default:
				data, err := s.netReader.ReadBytes('\n')
				if err != nil {
					log.Println("Net read error:", err)
					respError(s, []byte("Connection with server lost,\n\ttrying to reconnect...\n"))
					ioutil.ReadAll(s.conn)
					s.conn.Close()

					// trying to restablish connection
					conn, err := s.connect(5)
					if err != nil {
						respError(s, []byte("failed to reconnect to server,\n\ttry again later\n"))
						return
					}
					s.setConn(conn)
					respMsg(s, []byte("*** connection established"))

					// need to (re)auth
					err = s.SendAuth([]byte(s.login))
					if err != nil {
						respError(s, []byte(err.Error()))
					}
					continue
				}

				runCmd, args, err := cmds.Parse(data)
				if err != nil {
					log.Println("\nUnknown response type:", err)
					continue
				}

				err = runCmd(s, args)
				if err != nil {
					log.Println("\nWrong response:", err)
					continue
				}
			}
		}
	}()

	var run func(msg []byte) error
	for {
		run = s.SendMsg
		// need authorisation otherwise just send messages
		if len(s.login) == 0 {
			run = s.SendAuth
			fmt.Print("Enter your name: ")
		}

		text, err := s.cliReader.ReadBytes('\n')
		if err != nil {
			s.cliWriter.WriteString("unable to read from console")
			s.cliWriter.Flush()
			continue
		}
		if len(bytes.TrimSpace(text)) > 0 {
			// run the command
			err := run(text)
			if err != nil {
				respError(s, []byte(err.Error()))
			}
		} else {
			s.cliWriter.WriteString(s.bar())
			s.cliWriter.Flush()
		}
	}
}

// SendAuth authorize user by given login on server
func (s *Session) SendAuth(login []byte) error {
	err := s.SendCmd("AUTH", login)
	if err != nil {
		return err
	}
	s.login = string(bytes.TrimSpace(login))
	return nil
}

// SendMsg sends given msg into current channel
func (s *Session) SendMsg(msg []byte) error {
	err := s.SendCmd("SEND", msg)
	if err != nil {
		return err
	}
	s.cliWriter.WriteString(s.bar())
	s.cliWriter.Flush()
	return nil
}

// SendCmd send cmd command to server with given arg.
// arg can be nil (e.g. for QUIT)
func (s *Session) SendCmd(cmd string, arg []byte) error {
	s.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))

	s.netWriter.WriteString(cmd)

	arg = bytes.TrimSpace(arg)
	if len(arg) > 0 {
		s.netWriter.WriteByte(' ')
		s.netWriter.Write(arg)
	}

	s.netWriter.WriteByte('\n')
	err := s.netWriter.Flush()
	if err != nil {
		return fmt.Errorf("%s:%v", cmd, err)
	}
	s.lastCmd = cmd
	return nil
}

// Sutdown close the session and connection to server
func (s *Session) Sutdown() error {
	err := s.SendCmd("QUIT", nil)
	ioutil.ReadAll(s.conn)
	s.conn.Close()
	return err
}

func (s *Session) bar() string {
	return fmt.Sprintf("%s@%s > ", s.login, s.currentChannel)
}

func (s *Session) connect(retry int) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.Dial("tcp", s.sHost)
	if err != nil {
		for i := 0; i < retry; i++ {
			conn, err = dialer.Dial("tcp", s.sHost)
			if err == nil {
				return conn, nil
			}
			time.Sleep(time.Second * time.Duration(i) * 2)
		}
		return nil, fmt.Errorf("unable to connect to server: %v", err)
	}
	return conn, nil
}

func (s *Session) setConn(conn net.Conn) {
	s.conn = conn
	s.netReader = bufio.NewReader(conn)
	s.netWriter = bufio.NewWriter(conn)
}
