package main

import (
	"encoding/json"
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
	http.HandleFunc("/checkTicket", handleTicketCheck)

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "arabic-ocr-300518-e2c236268e78.json")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("server started at %s\n", ":"+port)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// check errors
func catch(err error) {
	if err != nil {
		fmt.Printf("error in %s\n", err)
		panic(err)
	}
}

func check(err error) {
	if err != nil {
		fmt.Printf("error in %s\n", err)
		panic(err)
	}
}

func handleTicketCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ticket check....")
	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	fmt.Println(id[0])
	job := GetJob(id[0])
	if job == nil {
		fmt.Println("ticket is invalid")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := &ClientUpdate{Status: 0}

	switch job.JobStatus {
	case 0: // waiting
		resp.Status = 0
	case 1: // running
		resp.Status = 1
	case 2: // done
		resp.Status = 2
		// serve fle to client here
	}
	json.NewEncoder(w).Encode(resp)
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
		fmt.Println("YOU DIDN'T UPLOAD A PDF, exiting")
		return
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	check(err)

	// Okay pdf file is uploaded. now return success to user and submit the job into redis queue
	freshJob := Job{JobStatus: 1, FileName: handler.Filename, FileData: fileBytes}
	jobID := GenRandomID() // uuid
	fmt.Println("submitting job")
	SubmitJob(jobID, freshJob)

	fmt.Println("sending ticket to client")
	ticket := &Ticket{
		Id:       jobID,
		FileName: handler.Filename,
	}

	json.NewEncoder(w).Encode(ticket)

	return

	//
	//	resultTextFileName := DoOCR(handler.Filename, fileBytes)
	//	fmt.Printf("resultTextFileName %s\n", resultTextFileName)
	//	// ....wait
	//	extension := filepath.Ext(resultTextFileName)
	//	name := resultTextFileName[0 : len(resultTextFileName)-len(extension)]
	//	serveFile(w, r, resultTextFileName)
	//	CleanDownloadedFiles(name)
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
}
