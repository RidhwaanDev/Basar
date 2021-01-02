package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	host     = "localhost"
	port     = "8000"
	OS_READ  = 04
	OS_WRITE = 02
)

func main() {
	StartServer()
}

func StartServer() {
	http.HandleFunc("/upload", handleUpload)
	fmt.Printf("server started at %s\n", host+":"+port)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
	defer CleanUpUploadsFolder()

	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range items {
		imageFile, err := os.Open("uploads/" + item.Name())
		check(err)
		img, err := jpeg.Decode(imageFile)
		check(err)

		fmt.Println("image decoded")
		jpeg.Encode(w, img, nil)

		fmt.Println("doing OCR")

		output, e := DetectText("uploads/" + item.Name())
		check(e)
		fmt.Println(output)
	}

}
