package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	err := ioutil.WriteFile("results/test1", []byte("hello how are you"), 0666)
	if err != nil {
		fmt.Println(err)
	}
}
