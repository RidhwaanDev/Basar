package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	host = "localhost"
	port = "8000"
)

func main() {
	StartServer()
}

func printenv() {
	envs := os.Environ()
	for _, a := range envs {
		str := strings.Split(a, "=")
		fmt.Printf("%s = %s\n", str[0], str[1])
	}
}

func StartServer() {
	http.HandleFunc("/upload", handleUpload)
	fmt.Printf("server started at %s\n", host+":"+port)
	// 	printenv()
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// check errors
func check(err error) {
	if err != nil {
		fmt.Printf("error in %s\n", err)
		panic(err)
	}
}

// uses uploads a pdf -> gets a pdf back in Arabic
func handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("error getting the file")
		fmt.Println(err)
		// panic(err)
		return
	}

	defer file.Close()
	// userUploadfileName := handler.Filename
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// check that they uploaded a pdf
	if filepath.Ext(handler.Filename) != ".pdf" {
		fmt.Println("YOU DIDN'T UPLOAD A PDF")
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	check(err)

	// create the file
	e := ioutil.WriteFile("pdf_to_convert.pdf", fileBytes, 0644)
	check(e)

	// convert pdf to a bunch of images and put them in the uploads directory
	// ConvertPDFToImages() puts a bunch of temp images in the uploads dir be sure to clean them out
	ConvertPDFToImages()
	// do the actual ocr
	Ocr()

	//	CleanUpUploadsFolder()

	fmt.Println("we are done")
}

func createFileToSendToClient(s string) *os.File {
	// create the file
	file, err := os.Create("input.txt")

	// write the OCR results to the file
	b := []byte(s)
	err = ioutil.WriteFile("input.txt", b, 0666)

	check(err)

	return file
}
