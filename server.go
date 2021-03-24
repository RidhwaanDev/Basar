package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	StartServer()
}

func StartServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", handleUpload)

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "arabic-ocr-300518-e2c236268e78.json")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("server started at %s\n", ":"+port)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// check errors
func check(err error) {
	if err != nil {
		fmt.Printf("error in %s\n", err)
		panic(err)
	}
}

func catch(err error) {
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

	resultTextFileName := DoOCR(handler.Filename, fileBytes)
	// ....wait
	serveFile(w, r, resultTextFileName)
}

func serveFile(writer http.ResponseWriter, request *http.Request, fileName string) {
	fmt.Println("func ServeFile in server.go")
	file, err := os.Open(fileName)
	check(err)
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader) // File offset is now len(fileHeader)
	fileType := http.DetectContentType(fileHeader)

	if err != nil {
		catch(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		catch(err)
	}

	writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s""`, fileInfo.Name()))
	writer.Header().Set("Content-Type", fileType)
	requestRange := request.Header.Get("range")
	fmt.Printf("request range: %s\n", requestRange)

	if requestRange == "" {
		writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
		file.Seek(0, 0)
		fmt.Println("filename path: ", filepath.Ext(fileInfo.Name()))
		io.Copy(writer, file)
	}
	// move filePath parsing into function
	extension := filepath.Ext(fileName)
	name := fileName[0 : len(fileName)-len(extension)]
	CleanDownloadedFiles(name)
}
