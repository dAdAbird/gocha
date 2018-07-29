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
		c.send([]byte(fmt.Sprintf("INFO *** %s is online\n", u.name)), nil)
		c.users[u] = struct{}{}
	}
	c.mx.Unlock()
}

func (c *Channel) UnRegister(u *user) {
	// send notification to other users
	c.mx.Lock()
	delete(c.users, u)
	c.mx.Unlock()

	c.SendString(fmt.Sprintf("INFO *** %s is offline\n", u.name), nil)
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

func (c *Channel) SendString(msg string, from *connection) {
	c.Send([]byte(msg), from)
}

func (c *Channel) Send(msg []byte, from *connection) {
	c.mx.Lock()
	c.send(msg, from)
	c.mx.Unlock()
}

func (c *Channel) send(msg []byte, from *connection) {
	for u := range c.users {
		for conn := range u.conns.c {
			// TODO: if error, close conn
			if conn != from {
				_, err := conn.Send(msg)
				if err != nil {

				}
			}
		}
	}
}
