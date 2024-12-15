package main

import (
	"fmt"
	"sync"
)

// 通过通信来共享内存，而不是通过共享内存来通信

type Map struct {
	m  map[int]int
	mu sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		m: make(map[int]int),
	}
}

func (m *Map) Set(key, val int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = val
}

func (m *Map) Get(key int) (int, bool) {
	m.mu.RLock()
	val, ok := m.m[key]
	m.mu.Unlock()
	return val, ok
}

func main() {
	m := NewGoodMap()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			m.setChan <- pair{key: key, val: key * 2}
		}(i)
	}

	wg.Wait()
	for i := 0; i < 10; i++ {
		respChan := make(chan resp)
		req := req{key: i, respChan: respChan}
		m.getChan <- req
		resp := <-respChan
		fmt.Println(resp.val, resp.ok)
	}
}

type pair struct {
	key int
	val int
}

type req struct {
	key      int
	respChan chan resp
}

type resp struct {
	val int
	ok  bool
}

type GoodMap struct {
	setChan chan pair
	getChan chan req
}

func NewGoodMap() *GoodMap {
	setChan := make(chan pair)
	getChan := make(chan req)

	go func() {
		m := make(map[int]int)

		for {
			select {
			case pair := <-setChan:
				m[pair.key] = pair.val
			case req := <-getChan:
				val, ok := m[req.key]
				req.respChan <- resp{val: val, ok: ok}
			}
		}
	}()

	return &GoodMap{
		setChan: setChan,
		getChan: getChan,
	}
}
