package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// remove everything from results dir
func main() {
	items, err := ioutil.ReadDir("results")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		// fmt.Println(item.Name())
		err := os.Remove("results/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}
}
