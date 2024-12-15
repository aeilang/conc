package main

import (
	"fmt"
	"sync"
)

// 通过共享内存来通信
func main1() {
	num := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			num++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(num)
}

// 通过通信来共享内存

type Num struct {
	addChan chan int
	getChan chan chan int
}

func NewNum() *Num {
	addChan := make(chan int)
	getChan := make(chan chan int)

	go func() {
		num := 0
		for {
			select {
			case val := <-addChan:
				num = num + val
			case respChan := <-getChan:
				respChan <- num
			}
		}
	}()

	return &Num{
		addChan: addChan,
		getChan: getChan,
	}
}

func (n *Num) Add(val int) {
	n.addChan <- val
}

func (n *Num) Get() int {
	respChan := make(chan int)
	n.getChan <- respChan
	val := <-respChan

	return val
}

func main() {
	num := NewNum()

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			num.Add(1)
		}()
	}

	wg.Wait()
	fmt.Println(num.Get())
}
