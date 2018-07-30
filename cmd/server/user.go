package main

import (
	"sync"
)

// User is chat user representation
type User struct {
	name  string
	conns cPool
}

// NewUser creates and returns new user by a given name
func NewUser(name string) *User {
	return &User{
		name: name,
		conns: cPool{
			c: make(map[*Client]struct{}),
		},
	}
}

// connections pool
// is a abstraction that keeps users' connections
type cPool struct {
	mx sync.Mutex
	c  map[*Client]struct{}
}

// Add adds a new connection to pool
func (p *cPool) Add(c *Client) {
	p.mx.Lock()
	p.c[c] = struct{}{}
	p.mx.Unlock()
}

// Len returns the pool's length
func (p *cPool) Len() int {
	p.mx.Lock()
	defer p.mx.Unlock()
	return len(p.c)
}

// Delete removes connection from pool and
// returns the remained connections count
func (p *cPool) Delete(c *Client) int {
	p.mx.Lock()
	defer p.mx.Unlock()
	delete(p.c, c)
	return len(p.c)
}
