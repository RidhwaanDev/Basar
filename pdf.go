package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	pdfToImagesCmdString = "cat *.go"
)

func main() {
	path, _ := convertPDFToImages("server.go")
	fmt.Println("path")
}

func isImageMagickInstalled() bool {
	path, err := exec.LookPath("convert")
	if err != nil {
		fmt.Println("error in looking for 'convert' cmd")
		log.Fatal(err)
		return false
	}

	log.Println(path)
	return true
}

func convertPDFToImages(path string) (string, error) {
	// check path

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// ocr-target is the last stop of the upload before it runs through OCR
	// ocr-target is a directory full of .pngs of the PDF

	err := os.Mkdir("ocr-target", os.ModeDir)

	//run the convert
	pdfToImagesCMD := exec.Command(pdfToImagesCmdString, "")

	var out bytes.Buffer

	pdfToImagesCMD.Stdout = &out

	err := pdfToImagesCMD.Run()

	if err != nil {
		log.Fatal(err)
	}

	return out.String()

}
