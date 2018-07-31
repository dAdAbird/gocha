package main

import (
	"fmt"
	"sync"
)

// Channel represent chat channels (something like Slack's channels)
type Channel struct {
	name  string
	mx    sync.Mutex
	users map[*User]struct{}
}

// NewChannel initialize and return new Channel
func NewChannel(name string) Channel {
	return Channel{
		name:  name,
		users: make(map[*User]struct{}),
	}
}

// Register adds user into channel
func (c *Channel) Register(u *User) {
	c.mx.Lock()
	if _, ok := c.users[u]; !ok {
		// send notification to other users
		c.send([]byte(fmt.Sprintf("INFO *** %s is online\n", u.Name)), nil)
		c.users[u] = struct{}{}
	}
	c.mx.Unlock()
}

// UnRegister removes user from channel
func (c *Channel) UnRegister(u *User) {
	c.mx.Lock()
	delete(c.users, u)
	c.mx.Unlock()

	// send notification to other users
	c.SendString(fmt.Sprintf("INFO *** %s is offline\n", u.Name), nil)
}

// Users returns list of users in chanel
func (c *Channel) Users() []*User {
	var users []*User
	c.mx.Lock()
	for u := range c.users {
		users = append(users, u)
	}
	c.mx.Unlock()
	return users
}

// SendString send message to the channel
func (c *Channel) SendString(msg string, from *Client) {
	c.Send([]byte(msg), from)
}

// Send send message to the channel
func (c *Channel) Send(msg []byte, from *Client) {
	c.mx.Lock()
	c.send(msg, from)
	c.mx.Unlock()
}

func (c *Channel) send(msg []byte, from *Client) {
	for u := range c.users {
		for _, conn := range u.Conns.Get() {
			// don't send message to sender
			if conn != from {
				// fmt.Printf("SEND: %s\n", msg)
				_, err := conn.Send(msg)
				if err != nil {
					// kick off the connection on error
					conn.colse <- fmt.Errorf("unable to send message: %v", err)
				}
			}
		}
	}
}
