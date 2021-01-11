package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type upload struct {
	name     string
	contact  string
	madrasah string
	filedata string
}

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

func handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server end point hit")

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	switch r.Method {

	case "GET":
		http.ServeFile(w, r, "index.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

// check errors
func check(err error) {
	if err != nil {
		fmt.Printf("error in %s\n", err)
		panic(err)
	}
}
