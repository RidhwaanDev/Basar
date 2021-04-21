package main

import (
	"bufio"
	"html/template"
	// "io"
	"net/http"
	"os"
)

type Todo struct {
	Txt string
}

type PageData struct {
	PageTitle string
	Lines     []Todo
}

// filename of the
func ServeResultAsHTML(filename string) {
	f, _ := os.Open(filename)
	reader := bufio.NewScanner(f)
	var lines []Todo
	for reader.Scan() {
		// fmt.Println(reader.Text())
		t := Todo{Txt: reader.Text()}
		lines = append(lines, t)
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			PageTitle: filename,
			Lines:     lines,
		}
		tmpl.Execute(w, data)
	})
}
