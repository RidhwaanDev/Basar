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
	// convert -density 150 input_file.pdf -quality 100 output_file.jpg
	prg := "convert"
	arg1 := "-density"
	val1 := "300"
	arg2 := "pdf_to_convert.pdf"
	arg3 := "-quality"
	val3 := "100"
	arg4 := ".jpg"

	cmd := exec.Command(prg, arg1, val1, arg2, arg3, val3, arg4)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// waits for command to complete
	err := cmd.Run()
	if err != nil {
		fmt.Println("an error in ConvertPDFToImages()" + "   " + stderr.String())
		log.Fatal(err)
	} else {
		// fmt.Println("moving converted images to upload dir")
		moveToUploadsDir()
	}
}

// move the jpgs to uploads directory
func moveToUploadsDir() {
	items, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println(err)
	}
	// this needs to be more foolproof
	for _, item := range items {
		if filepath.Ext(item.Name()) == ".jpg" {
			os.Rename(item.Name(), "uploads/"+item.Name())
		}
	}
}

func ConvertTextToPDF(s string) {
	// chrome --headless --disable-gpu --print-to-pdf text.txt

	// convert the string to bytes
	b := []byte(s)

	err := ioutil.WriteFile("text.txt", b, 0644)
	check(err)

	prg := "chrome"
	arg1 := "--headless"
	arg2 := "--disable-gpu"
	arg3 := "--print-to-pdf"
	val1 := "text.txt"
	cmd := exec.Command(prg, arg1, arg2, arg3, val1)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("an error in pdf convert" + "   " + stderr.String())
		log.Fatal(err)
	} else {
		fmt.Println("successfully converted .txt to .pdf")
	}
}
