package main

import (
	"sync"
)

// User is chat user representation
type User struct {
	Name  string
	Conns cPool
}

type users struct {
	mx   sync.Mutex
	pool map[string]*User
}

// global users pool
var connectedUsers users

func init() {
	connectedUsers.pool = make(map[string]*User)
}

// UserAuth authorise user on server
func UserAuth(name string, c *Client) *User {
	return connectedUsers.addOrGet(name, c)
}

// UserLogout removes user from server
func UserLogout(u *User) {
	connectedUsers.remove(u.Name)
}

func (us *users) addOrGet(name string, c *Client) *User {
	us.mx.Lock()
	defer us.mx.Unlock()

	u, ok := us.pool[name]
	if !ok {
		u = &User{
			Name: name,
			Conns: cPool{
				c: make(map[*Client]struct{}),
			},
		}
	}
	u.Conns.Add(c)
	us.pool[name] = u

	return us.pool[name]
}

func (us *users) remove(name string) {
	us.mx.Lock()
	delete(us.pool, name)
	us.mx.Unlock()
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

// Get returns slice of available connections
func (p *cPool) Get() []*Client {
	var clients []*Client
	p.mx.Lock()
	for u := range p.c {
		clients = append(clients, u)
	}
	p.mx.Unlock()
	return clients
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
