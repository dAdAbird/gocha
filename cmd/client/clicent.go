package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
)

func main() {
	lhost := "localhost:1313"

	dialer := new(net.Dialer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", lhost)
	if err != nil {
		log.Fatalln("Unable to start listener:", err)
	}
	defer conn.Close()

	cmd := make(CommadsSet)
	cmd.Register("MSG", RespMsg)
	cmd.Register("INFO", RespMsg)
	cmd.Register("OK", RespOK)
	cmd.Register("ERR", RespError)

	cn := &connection{
		conn:      conn,
		pipe:      make(chan struct{}, 1),
		netReader: bufio.NewReader(conn),
		netWriter: bufio.NewWriter(conn),
		cliReader: bufio.NewReader(os.Stdin),
		cliWriter: bufio.NewWriter(os.Stdout),
	}

	cn.Process(&cmd)

}
