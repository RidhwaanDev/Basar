package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
	userUploadfileName := handler.Filename
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

		imageFile, err := os.Open("uploads/" + item.Name())
		check(err)
		// image assumed to be a .jpg. if its not, this will break
		img, err := jpeg.Decode(imageFile)
		check(err)

		fmt.Println("image decoded")
		jpeg.Encode(w, img, nil)

		fmt.Println("doing OCR")
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(1 * time.Second)
			resc <- "hi"
		}()
		// go DetectText("uploads/"+item.Name(), &wg, resc)
	}

	// wait for goroutines to finish
	go func() {
		wg.Wait()
		close(resc)
	}()

	var b strings.Builder
	b.Grow(100)

	for str := range resc {
		b.WriteString(str)
	}
	// put OCR results in a .txt file, returns *os.File and its length
	finalFile, _, size := createFileToSendToClient(b.String())

	//prepare to send .txt file to client

	nameWithoutExt := strings.Split(userUploadfileName, ".")[0]

	attachment := fmt.Sprintf("attachment; filename=%s.txt", nameWithoutExt)
	fmt.Println(attachment)

	w.Header().Set("Content-Disposition", "attachment; filename=YourFile")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	check(err)
	fmt.Printf("finalFile size: %d\t bytes written: %d\n", size, n)
	serveFile(w, r)
}

func createFileToSendToClient(s string) (*os.File, []byte, int64) {
	// create the file
	file, err := os.Create("finalres.txt")

	// write the OCR results to the file
	b := []byte(s)
	err = ioutil.WriteFile("finalres.txt", b, 0644)
	check(err)

	fmt.Println("wrote the final .txt file, nice!")

	f, err := file.Stat()
	size := f.Size()
	check(err)

	// return the *os.File and its length

	return file, b, size
}
