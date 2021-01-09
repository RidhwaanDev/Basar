package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	host     = "localhost"
	port     = "8000"
	OS_READ  = 04
	OS_WRITE = 02
	TXT      = 1
	PDF      = 2
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "txt":
			format = TXT
		case "pdf":
			format = PDF
		case "png":
			format = PNG
		case "jpg":
			format = JPG
		default:
			format = PDF
		}
	}

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

// the only file type supported (as of now) is .pdf
// user should upload a single .pdf -> convert into images -> do ocr -> send .pdf bacl
func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("file upload endpoint hit")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("error getting the file")
		fmt.Println(err)
		panic(err)
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
	ConvertPDFToImages()

	// ConvertPDFToImages() puts a bunch of temp images in uploads, dir be sure to clean them out
	defer CleanUpUploadsFolder()

	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}
	// for each image, spawn a gorotuine to do OCR on it. Manage all these gorutines with
	// a waitgroup
	var wg sync.WaitGroup
	// channel of all the OCR results.
	resc := make(chan string)
	for _, item := range items {
		fmt.Println("doing OCR")
		wg.Add(1)
		go DetectText("uploads/"+item.Name(), &wg, resc)
	}

	// wait for goroutines to finish
	go func() {
		wg.Wait()
		close(resc)
	}()

	// string builder for all the OCR results
	var b strings.Builder
	b.Grow(100)

	// iter chan of OCR results
	for str := range resc {
		b.WriteString(str)
	}

	// put OCR results in a .txt file and return the *os.File object
	clientFileTxt := createFileToSendToClient(b.String())
	// if we want to send client a .pdf file
	if format == PDF {
		// conver the  .txt file into a pdf, and get the pdf file, this uses headless chrome , see pdf.go
		convertedPDFFile := ConvertTextToPDF(clientFileTxt)
		// serve the file to the client
		ServeFile(w, r, convertedPDFFile)
	} else {
		// .txt
		ServeFile(w, r, clientFileTxt)
	}

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
