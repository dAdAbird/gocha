package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	GitCommit string
	BuildTime string
	Version   string

	listenOn    = flag.String("h", ":1313", "Chat server host")
	showVersion = flag.Bool("version", false, "show version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nBuild time: %s\nCommit: %s\n", Version, BuildTime, GitCommit)
		os.Exit(0)
	}

	ln, err := net.Listen("tcp", *listenOn)
	if err != nil {
		log.Fatalln("Unable to start listener:", err)
	}
	defer ln.Close()
	log.Println("Listening on", *listenOn)

	// #General channel
	defaulChannel := NewChannel("general")

	cmd := make(CommadsSet)
	cmd.Register("AUTH", cmdAuth)
	cmd.Register("SEND", cmdWrapAuth(cmdSend))
	cmd.Register("QUIT", cmdQuit)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error while connection accept:", err)
			continue
		}

		go func(c net.Conn) {
			cn := Client{
				conn:        c,
				state:       scInit,
				currChannel: &defaulChannel,
				colse:       make(chan error, 1),
			}
			cn.Handle(&cmd)
		}(conn)
	}
}
