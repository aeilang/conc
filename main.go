package main

import (
	"fmt"
	"sync"
)

// 通过通信来共享内存，而不是通过共享内存来通信
func main() {
	num := 0

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			num++
		}()
	}

	wg.Done()
	fmt.Println(num)
}
