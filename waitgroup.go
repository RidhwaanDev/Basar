package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	resc := make(chan string)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("doing work... %d\n", id)
			time.Sleep(2 * time.Second)
			resc <- generateRandomString() + "$" + string(id)
		}(i)

	}

	go func() {
		wg.Wait()
		close(resc)
	}()

	for str := range resc {
		fmt.Println(str)
	}
}

func generateRandomString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
