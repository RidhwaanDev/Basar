package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// converts the pdf to a series of images. puts those images in the uploads dir
func ConvertPDFToImages() {
	prg := "convert"
	arg1 := "-density"
	val1 := "150"
	arg2 := "pdf_to_convert.pdf"
	arg3 := "-quality"
	val3 := "100"
	arg4 := "output_file.jpg"

	cmd := exec.Command(prg, arg1, val1, arg2, arg3, val3, arg4)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("an error in convert" + "   " + stderr.String())
		log.Fatal(err)
	} else {
		moveToUploadsDir()
	}
}

// move the jpgs to uploads directory
func moveToUploadsDir() {
	items, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range items {
		if filepath.Ext(item.Name()) == ".jpg" {
			// fmt.Println(item.Name())
			os.Rename(item.Name(), "uploads/"+item.Name())
		}
	}
}
