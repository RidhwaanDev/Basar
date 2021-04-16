package main

import (
	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	StartServer()
}

func StartServer() {
	root := "/home/ridhwaan/ArabicOCR/static/"
	http.Handle("/", http.FileServer(http.Dir(root)))
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/checkTicket", handleTicketCheck)

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "arabic-ocr-300518-e2c236268e78.json")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("server started at %s\n", ":"+port)

	dir, _ := os.Getwd()
	fmt.Println("current path :" + dir)

	log.Fatal(http.ListenAndServe(":8000", nil))

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

func handleJobFileRemoval(w http.ResponseWriter, r *http.Request) {
	// client  downloaded file, time to remove the text file
	enableCors(&w)
	fmt.Println("hit handle jobFileRemoval")

	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	err := os.Remove(fmt.Sprintf("./static/%s.txt", id[0]))
	check(err)
}

func handleTicketCheck(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println("hit handleTicetCheck")

	id, ok := r.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	job := GetJob(id[0])
	if job == nil {
		fmt.Println("ticket is invalid")
		return
	}
	// fmt.Printf("ticket check with id: %s", id[0])

	// w.Header().Set("Content-Type", "application/json")
	resp := &ClientUpdate{Status: 0, FileName: job.FileName}

	switch job.JobStatus {
	case 0: // waiting
		resp.Status = 0
		json.NewEncoder(w).Encode(resp)
	case 1: // running
		resp.Status = 1
		json.NewEncoder(w).Encode(resp)
	case 2: // done
		resp.Status = 2
		// serve fle to client here
		fmt.Println("OCR is done, this is it boys, send it!")
		// http.ServeFile(w, r, job.FileName)
		// json.NewEncoder(w).Encode(resp)

		file := job.FileName
		fmt.Printf("just printing the file: %s \n", file)
		fileNewLocation := fmt.Sprintf("./static/%s", file)
		err := os.Rename(file, fileNewLocation)
		check(err)

		// tell the client the job is done. They shoulkd be able to to read the file from ./static
		json.NewEncoder(w).Encode(resp)

		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	}
}

// uses uploads a pdf -> gets a pdf back in Arabic
func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit download endpoint")
	enableCors(&w)

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
	// TODO limit fileSize
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// check that they uploaded a pdf
	if filepath.Ext(handler.Filename) != ".pdf" {
		fmt.Println("YOU DIDN'T UPLOAD A PDF, exiting")
		return
	}

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	check(err)

	//send ticket to user with job ID and submit the job into redis
	jobID := GenRandomID() // uuid
	freshJob := Job{JobStatus: 1, FileName: jobID + ".txt", FileData: fileBytes}
	fmt.Println("submitting job")
	SubmitJob(jobID, freshJob)

	fmt.Println("sending ticket to client")
	// ticket for the user
	ticket := &Ticket{
		Id:       jobID,
		FileName: handler.Filename,
	}

	json.NewEncoder(w).Encode(ticket)

	// go do the OCR and return
	go DoOCR(jobID, handler.Filename, fileBytes)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
