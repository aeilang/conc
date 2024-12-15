package main

import (
	"fmt"
	"sync"
)

// 通过通信来共享内存，而不是通过共享内存来通信
func main() {
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
