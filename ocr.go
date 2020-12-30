package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

var repeat_map = make(map[string]bool)

const (
	DESKTOP_PATH = "/Users/ridhwaananayetullah/Desktop/"
)

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("send location")
		return
	}
	path := DESKTOP_PATH + os.Args[1]

	dir, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening directory: %s", err)
	}
	defer dir.Close()

	list, _ := dir.Readdirnames(0)
	res := make([]string, 100)
	resc := make(chan string)
	for _, name := range list {
		var filePath = path + "/" + name
		go DetectText(filePath, resc)
	}

	close(resc)

	// collect results from OCR goroutine
	for {
		res = append(res, <-resc)
	}

	fmt.Printf("final result: %s\n", res)

}

func DetectText(file string, resc chan<- string) error {
	fmt.Printf("detecting text in %s\n", file)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		log.Fatal(err)
		return err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("No text found.")
	outputString := make([]string, 50)
	cnt := 0
	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		// fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			cnt++
			repeat_map[annotation.Description] = true
			// the first line is the ocr of the entire document
			outputString = append(outputString, annotation.Description)
			break
			// fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	output := strings.Join(outputString, "\n")
	resc <- output
	return nil
}
