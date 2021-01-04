package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(2 * time.Second)
			fmt.Printf("doing work %d\n", id)
		}(i)
	}
	wg.Wait()
}
