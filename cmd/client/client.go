package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	gitCommit string
	buildTime string
	version   string

	serverHost  = flag.String("h", "localhost:1313", "Chat server host")
	showVersion = flag.Bool("version", false, "show version and exit")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nBuild time: %s\nCommit: %s\n", version, buildTime, gitCommit)
		os.Exit(0)
	}

	// create session
	se := &Session{
		sHost:     *serverHost,
		close:     make(chan struct{}, 1),
		cliReader: bufio.NewReader(os.Stdin),
		cliWriter: bufio.NewWriter(os.Stdout),
	}

	conn, err := se.connect(5)
	if err != nil {
		log.Fatalln("Unable to connect to server:", err)
	}
	defer conn.Close()
	se.setConn(conn)

	se.cliWriter.WriteString("Connected to server\n")
	se.cliWriter.Flush()

	signalsHandler(se)

	rsps := make(ResponsesSet)
	rsps.Register("MSG", respMsg)
	rsps.Register("INFO", respMsg)
	rsps.Register("OK", respOK)
	rsps.Register("BYE!", respOK)
	rsps.Register("ERR", respError)

	fmt.Println("To exit, press '^C'")
	fmt.Println()

	se.Process(&rsps)
}

// System signal handling
func signalsHandler(s *Session) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-c
		fmt.Printf("\nSignal %v was received, closing connection to server\n", sig)
		s.close <- struct{}{}
		fmt.Println("BYE!")
		os.Exit(0)
	}()
}
