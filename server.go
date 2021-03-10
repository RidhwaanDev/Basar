package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	host = "localhost"
	port = "8000"
)

func main() {
	StartServer()
}

func StartServer() {
	http.HandleFunc("/upload", handleUpload)
	fmt.Printf("server started at %s\n", host+":"+port)
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
	defer CleanUpUploadsFolder()

	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}
	// for each image, spawn a gorotuine to do OCR on it. Manage all these goroutines with
	// a waitgroup
	var wg sync.WaitGroup
	// channel of all the OCR results.
	resc := make(chan string)
	for _, item := range items {
		// fmt.Printf("doing OCR on item: %s\n", item.Name())
		wg.Add(1)
		go DetectText("uploads/"+item.Name(), &wg, resc)
	}

	//	defer CleanUpResultsFolder()
	// wait for goroutines to finish
	go func() {
		wg.Wait()
		close(resc)
	}()

	finalOutputTextFile := CreateFinalOutputTextFile()

	//serve the file to the he user
	ServeFile(w, r, finalOutputTextFile)
	//	// put OCR results in a .txt file and return the *os.File object
	//	clientFileTxt := createFileToSendToClient(b.String())
	//	// if we want to send client a .pdf file
	//	if format == PDF {
	//		// conver the  .txt file into a pdf, and get the pdf file, this uses headless chrome , see pdf.go
	//		convertedPDFFile := ConvertTextToPDF(clientFileTxt)
	//		// serve the file to the client
	//		ServeFile(w, r, convertedPDFFile)
	//	} else {
	//		// .txt
	//		ServeFile(w, r, clientFileTxt)
	//	}
	// we are done :)
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
