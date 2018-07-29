package main

import (
	"log"
	"net"
)

func main() {
	lhost := ":1313"
	ln, err := net.Listen("tcp", lhost)
	if err != nil {
		log.Fatalln("Unable to start listener:", err)
	}
	defer ln.Close()
	log.Println("Listening on", lhost)

	// #General channel
	defaulChannel := NewChannel("general")

	cmd := make(CommadsSet)
	cmd.Register("AUTH", CmdAuth)
	cmd.Register("SEND", cmdWrapAuth(CmdSend))
	cmd.Register("QUIT", CmdQuit)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error while connection accept:", err)
			continue
		}

		go func(c net.Conn) {
			cn := connection{
				c:           c,
				state:       scInit,
				currChannel: &defaulChannel,
				colse:       make(chan struct{}, 1),
			}
			// err := cn.c.SetWriteDeadline(30 * time.Second)
			// if err != nil {
			// 	log.Printf("Unable to set Write Deadline for conn from %v: %v", c.RemoteAddr(), err)
			// }
			cn.Process(&cmd)
		}(conn)
	}
}
