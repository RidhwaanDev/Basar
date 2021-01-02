package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

func DetectText(file string) (string, error) {
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

	fmt.Println("No text found.")
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
			break
			// fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	output := strings.Join(outputString, "\n")
	return output, nil
}
