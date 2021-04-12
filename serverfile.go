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

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.

	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

	//stream the body to the client without fully loading it into memory
	io.Copy(w, f)
}
