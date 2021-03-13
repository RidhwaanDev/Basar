package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	//"strings"
	"sync"
)

// for each image, spawn a gorotuine to do OCR on it. Manage all these goroutines with
func Ocr() {
	items, err := ioutil.ReadDir("uploads")
	check(err)
	// a waitgroup
	var wg sync.WaitGroup
	// channel of all the OCR results.
	resc := make(chan map[string]string)
	for _, item := range items {
		wg.Add(1)
		// fmt.Println(item.Name())
		go DetectText("uploads/"+item.Name(), &wg, resc)
	}
	//	defer CleanUpResultsFolder()
	// wait for goroutines to finish
	go func() {
		wg.Wait()
		close(resc)
	}()

	for str := range resc {
		for key, val := range str {
			fmt.Println(key + " " + val)
		}
	}
}

func DetectText(fileName string, wg *sync.WaitGroup, resc chan<- map[string]string) (string, error) {
	defer wg.Done()

	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println("Opening file: " + fileName)
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer f.Close()

	// don't use a file. maybe keep an io.Reader in memory somewhere?
	// TODO check performance diff in memory vs reading from a file
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		check(err)
		log.Fatal(err)
		return "", err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		check(err)
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
			// TODO useful for line by line transcription
			// the first line is the ocr of the entire document
			outputString = append(outputString, annotation.Description)
			// fmt.Println(outputString)
			// resc <- strings.Join(outputString, "\n")
			// resc <- outputString
			fileNameToText := make(map[string]string)
			fileNameToText[fileName] = annotation.Description
			resc <- fileNameToText
			break
		}
	}

	// writeToResult(fileName, output)

	return "test", nil
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
