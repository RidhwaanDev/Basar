package main

import (
	"os/exec"
	"fmt"
)

type pdf struct {
	location string // location on disk
	name     string
	pages    int
	size     int64
}

// number of pages in the PDF
func (pdf p) GetPageCount(ch <-chan int) {}

// mb size of PDF, usually when client uploads pdf, we will already know that
func (pdf p) GetSize(ch <-chan int64) {}

// headless chrome lol
func ConvertTextToPDF {
	// chrome --headless --disable-gpu --print-to-pdf input.txt
	prg := "chrome"
	arg1 := "--headless"
	arg2 := "--disable-gpu"
	arg3 := "--print-to-pdf"
	val1 := "final_res.txt"
	cmd := exec.Command(prg, arg1, arg2, arg3, val1)

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("an error in pdf convert" + "   " + stderr.String())
		log.Fatal(err)
	} else {
		fmt.Println("successfully converted .txt to .pdf")
	}
}
