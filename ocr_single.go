package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("need full file name")
		return
	}

	ocr(os.Args[1])
}

func ocr_dir(file string) {

	if len(os.Args) <= 1 {
		fmt.Println("send location")
		return
	}
	path := "shaa"

	dir, err := os.Open(path)

	if err != nil {
		log.Fatalf("Error opening directory: %s", err)
	}

	defer dir.Close()
	list, _ := dir.Readdirnames(0)
	for _, name := range list {
		var filePath = path + "/" + name
		ocr(filePath)
	}
}

func ocr(file string) {
	fmt.Printf("detecting text in %s\n", file)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("No text found.")
	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		// fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			// the first line is the ocr of the entire document
			fmt.Println(annotation.Description)
			break
			// fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}
}
