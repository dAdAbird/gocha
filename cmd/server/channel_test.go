package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"testing"
)

func TestChannelAdd(t *testing.T) {
	chn := NewChannel("test")

	addCnt := 20

	var wg sync.WaitGroup
	for i := 0; i < addCnt; i++ {
		wg.Add(1)
		go func() {
			chn.Register(UserAuth(fmt.Sprintf("user%d", rand.Intn(1e3)), &Client{conn: &net.TCPConn{}, colse: make(chan error, 3e2)}))
			wg.Done()
		}()
		// check for races
		u := UserAuth(fmt.Sprintf("user_%d", rand.Intn(1e3)), &Client{conn: &net.TCPConn{}, colse: make(chan error, 3e2)})
		chn.Register(u)
		chn.UnRegister(u)
	}
	wg.Wait()

	users := chn.Users()
	if len(users) != addCnt {
		t.Errorf("Wrong users count, got %d has to be %d", len(users), addCnt)
	}
}

func TestChannelBroadcast(t *testing.T) {

	result := `INFO *** testuser1 is online
MSG testuser1:	Hello
MSG testuser0:	Hello
MSG testuser1:	Hello
MSG testuser0:	Hello
`

	l, err := net.Listen("tcp", ":13131")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	// var cn *Client
	var cns [5]*Client
	chn := NewChannel("test")
	go func() {
		for i := 0; i < 5; i++ {
			conn, err := l.Accept()
			if err != nil {
				t.Fatal("unable to accept", err)
			}
			defer conn.Close()

			cns[i] = &Client{
				conn:        conn,
				state:       scInit,
				currChannel: &chn,
			}

			_, err = cmdAuth(cns[i], []byte(fmt.Sprintf("testuser%d", i%2)))
			if err != nil {
				t.Fatal("unable to AUTH", err)
			}

			_, err = cmdSend(cns[i], []byte("Hello"))
			if err != nil {
				t.Fatal("unable to SEND msg", err)
			}
		}
	}()

	var connc [5]net.Conn
	for i := 0; i < 5; i++ {
		connc[i], err = net.Dial("tcp", ":13131")
		if err != nil {
			t.Fatal("unable to dial", err)
		}
		defer connc[i].Close()
	}

	// var msg []byte
	// _, err = connc.Read(msg)
	msg, err := ioutil.ReadAll(connc[0])
	if err != nil {
		t.Fatal("unable to read stream", err)
	}

	if string(msg) != result {
		t.Errorf("Wrong data in stream got:\n%s\nhas to be\n%s", msg, result)
	}
}
