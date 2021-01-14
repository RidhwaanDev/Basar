package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func DetectText(file string, wg *sync.WaitGroup, resc chan<- string) (string, error) {
	defer wg.Done()
	fmt.Printf("detecting text in %s\n", file)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	outputString := make([]string, 50)
	cnt := 0
	if len(annotations) == 0 {
		fmt.Println("No text found.")
	} else {
		// fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			cnt++
			// the first line is the ocr of the entire document
			outputString = append(outputString, annotation.Description)
			// fmt.Println(outputString)
			resc <- strings.Join(outputString, "\n")
			break
		}
	}

	output := strings.Join(outputString, "\n")

	writeToResult(file, output)

	return output, nil
}

// write each OCR result in its own file and put it in into the results directory
func writeToResult(filename string, result string) {
	_, fname := filepath.Split(filename)
	fmt.Println(fname)
	err := ioutil.WriteFile("results/"+fname+".txt", []byte(result), 0666)
	if err != nil {
		fmt.Printf("error in writeToResult %s\n", err)
	}
}
