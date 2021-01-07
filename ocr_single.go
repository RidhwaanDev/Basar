package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// sync file i/o
var mutex = &sync.Mutex{}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("please provide a file or directory name")
		return
	}

	name := os.Args[1]
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		ocr_dir(name)
	case mode.IsRegular():
		// single file
		//		detectText(name)
	}
}

func ocr_dir(file string) {
	path := "images"

	dir, err := os.Open(path)

	if err != nil {
		log.Fatalf("Error opening directory: %s", err)
	}

	f, err := os.Create("output.txt")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer dir.Close()
	list, _ := dir.Readdirnames(0)
	for _, name := range list {
		var filePath = path + "/" + name
		go detectText(filePath, f)
	}
}

// writes to output.txt
func detectText(file string, f *os.File) {
	fmt.Printf("detecting text in %s\n", file)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatal(err)
	}

	outputString := make([]string, 50)

	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		// fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			// the first line is the ocr of the entire document
			outputString = append(outputString, annotation.Description)
			// fmt.Println(annotation.Description)
			break
			// fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	output := strings.Join(outputString, " ")

	// sync file i/o
	mutex.Lock()
	// write the OCR results to the file
	if _, err := f.WriteString(output); err != nil {
		log.Println(err)
	}
	mutex.Unlock()

	if err != nil {
		fmt.Println(err)
	}
}
