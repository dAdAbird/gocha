package main

import (
	"errors"
	"sync"
)

type user struct {
	id    int
	name  string
	conns cPool
}

var ErrUserNoName = errors.New("user name doesn't provided")

type users struct {
	mx   sync.Mutex
	pool map[string]*user
}

var connectedUsers users

func init() {
	connectedUsers.pool = make(map[string]*user)
}

func (us *users) Remove(name string) {
	us.mx.Lock()
	delete(us.pool, name)
	us.mx.Unlock()
}

func (us *users) Get(u *user) {

}

func (us *users) Add(name string, c *connection) *user {
	us.mx.Lock()
	defer us.mx.Unlock()

	u, ok := us.pool[name]
	if !ok {
		u = &user{
			name: name,
			conns: cPool{
				c: make(map[*connection]struct{}),
			},
		}
	}
	u.conns.Add(c)
	us.pool[name] = u

	return us.pool[name]
}
