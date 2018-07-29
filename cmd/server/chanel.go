package main

import (
	"fmt"
	"sync"
)

type Channel struct {
	name  string
	mx    sync.Mutex
	users map[*user]struct{}
}

func NewChannel(name string) Channel {
	return Channel{
		name:  name,
		users: make(map[*user]struct{}),
	}
}

func (c *Channel) Register(u *user) {
	c.mx.Lock()
	// send notification to other users
	if _, ok := c.users[u]; !ok {
		c.send([]byte(fmt.Sprintf("*** %s is online\n", u.name)))
		c.users[u] = struct{}{}
	}
	c.mx.Unlock()
}

func (c *Channel) UnRegister(u *user) {
	// send notification to other users
	c.mx.Lock()
	delete(c.users, u)
	c.mx.Unlock()

	c.SendString(fmt.Sprintf("*** %s is offline\n", u.name))
}

func (c *Channel) Users() []*user {
	var users []*user
	c.mx.Lock()
	for u := range c.users {
		users = append(users, u)
	}
	c.mx.Unlock()
	return users
}

func (c *Channel) SendString(msg string) {
	c.Send([]byte(msg))
}

func (c *Channel) Send(msg []byte) {
	c.mx.Lock()
	c.send(msg)
	c.mx.Unlock()
}

func (c *Channel) send(msg []byte) {
	for u := range c.users {
		for conn := range u.conns.c {
			// TODO: if error, close conn
			// if c != conn {
			_, err := conn.Send(msg)
			if err != nil {

			}
			// }
		}
	}
}
