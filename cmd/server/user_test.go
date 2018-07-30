package main

import (
	"sync"
	"testing"
)

func TestConnPool(t *testing.T) {
	pool := &cPool{
		c: make(map[*Client]struct{}),
	}
	cl := new(Client)
	addCnt := 20

	var wg sync.WaitGroup
	for i := 0; i < addCnt; i++ {
		wg.Add(1)
		go func() {
			pool.Add(&Client{})
			wg.Done()
		}()
		// check for races
		pool.Add(cl)
		pool.Delete(cl)
		pool.Len()
	}
	wg.Wait()

	if pool.Len() != addCnt {
		t.Errorf("Wrong pool length, got %d has to be %d", pool.Len(), addCnt)
	}
}
