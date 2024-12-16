package main

import (
	"context"
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
	getChan chan int
	cancel  context.CancelFunc
}

func NewNum() *Num {
	addChan := make(chan int)
	getChan := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		num := 0
		for {
			select {
			case val, ok := <-addChan:
				if !ok {
					return
				}
				num = num + val
			case getChan <- num:
			case <-ctx.Done():
				close(addChan)
				close(getChan)
				return
			}
		}
	}()

	return &Num{
		addChan: addChan,
		getChan: getChan,
		cancel:  cancel,
	}
}

func (n *Num) Add(val int) {
	n.addChan <- val
}

func (n *Num) Get() int {
	val := <-n.getChan

	return val
}

// Close 关闭Num，清理资源
func (n *Num) Close() {
	n.cancel()
}

func main() {
	num := NewNum()
	defer num.Close()

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
