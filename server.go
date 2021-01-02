package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

	tempFile, err := ioutil.TempFile("uploads", "upload-*.jpg")
	// clear out the uploads directory
	defer CleanUpTemp()

	defer tempFile.Close()

	check(err)

	// get the bytes of the file the user uploaded
	fileBytes, err := ioutil.ReadAll(file)
	check(err)

	n, err := tempFile.Write(fileBytes)
	check(err)

	fmt.Printf("wrote %d bytes to tempFile!\n", n)

	fileInfo, err := tempFile.Stat()
	check(err)
	// get the file extension. ".png"
	fileType := strings.Split(fileInfo.Name(), ".")[1]
	fmt.Printf("file type for %s : %s\n", tempFile.Name(), fileType)

	if fileType == "png" || fileType == "jpg" {
		fmt.Println("found image\n")

		imageFile, err := os.Open(tempFile.Name())
		check(err)
		img, err := jpeg.Decode(imageFile)
		check(err)

		fmt.Println("image decoded")
		jpeg.Encode(w, img, nil)

		fmt.Println("doing OCR")

		output, e := DetectText(tempFile.Name())
		check(e)
		fmt.Println(output)
	}
}
