package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	// "os"
	"path/filepath"
	"strings"
)

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func ServeFile(w http.ResponseWriter, r *http.Request, name string) {
	file := fileNameWithoutExtension(name)
	file = file + ".txt"
	fmt.Printf("just printing the file: %s \n", file)
	f, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	stat, _ := f.Stat()

	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", stat.Size())

	//stream the body to the client without fully loading it into memory
	io.Copy(w, f)
}

// func main() {
// 	http.HandleFunc("/", Home)
// 	http.HandleFunc("/download", ForceDownload)
//
// 	// SECURITY : Only expose the file permitted for download.
// 	http.Handle("/"+file, http.FileServer(http.Dir("./")))
//
// 	http.ListenAndServe(":8080", nil)
// }
