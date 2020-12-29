package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CleanUpTemp() {
	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		fmt.Println(item.Name())
		err := os.Remove("uploads/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}
}