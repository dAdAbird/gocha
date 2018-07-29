package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type connection struct {
	userName       string
	currentChannel string
	conn           net.Conn
	pipe           chan struct{}
	lastCmd        string
	netReader      *bufio.Reader
	netWriter      *bufio.Writer
	cliReader      *bufio.Reader
	cliWriter      *bufio.Writer
}

func (c *connection) Process(cmds *CommadsSet) {
	go func() {
		for {
			data, err := c.netReader.ReadBytes('\n')
			if err != nil {
				log.Println("Net read error:", err)
			}

			runCmd, args, err := cmds.Parse(data)
			if err != nil {
				log.Println("Unknown responce type:", err)
				continue
			}

			_, err = runCmd(c, args)
			if err != nil {
				log.Println("Wrong responce:", err)
				continue
			}
		}
	}()

	for {
		// select {
		// case <-c.pipe:
		// 	c.cliWriter.Flush()
		// default:
		if len(c.userName) == 0 {
			fmt.Print("Type your name: ")
			text, _ := c.cliReader.ReadString('\n')
			c.netWriter.WriteString("AUTH " + text)
			c.netWriter.Flush()
			c.lastCmd = "AUTH"
			c.userName = strings.TrimSpace(text)
		} else {
			// fmt.Printf("%s@%s > ", c.userName, c.currentChannel)
			text, _ := c.cliReader.ReadString('\n')
			if len(strings.TrimSpace(text)) > 0 {
				c.netWriter.WriteString("SEND " + text)
				c.netWriter.Flush()
				c.lastCmd = "SEND"
			}
		}
		// }
	}
}
