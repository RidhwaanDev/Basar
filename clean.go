package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func CleanUpUploadsFolder() {
	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		// fmt.Println(item.Name())
		err := os.Remove("uploads/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}
	// remove the pdf file
	os.Remove("pdf_to_convert.pdf")
}

// remove everything from results dir
func CleanUpResultsFolder() {
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
