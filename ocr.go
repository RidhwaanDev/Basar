package main

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"io"
	"os"
)

var repeat_map = make(map[string]bool)

//func main() {
//	// detectText(os.Stdout, "MunirahImanPage6.png")
//	// file, err := os.Open("images/iman")
//	//	if err != nil {
//	//		log.Fatalf("Error opening directory: %s", err)
//	//	}
//	//	defer file.Close()
//	//
//	//	list, _ := file.Readdirnames(0)
//	//	for _, name := range list {
//	//		fmt.Println(name)
//	//	}
//}

func Print() {
	fmt.Println("hello")
}
func DetectText(w io.Writer, file string, errors chan<- error) error {
	fmt.Printf("detecting text in %s\n", file)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		errors <- err
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		errors <- err
		return err
	}

	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		errors <- err
		return err
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		errors <- err
		return err
	}

	outputString := make([]string, 50)
	cnt := 0
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		// fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			cnt++
			repeat_map[annotation.Description] = true
			outputString = append(outputString, annotation.Description)
			break
			// fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	fmt.Println(outputString)
	fmt.Println(cnt)
	errors <- nil
	return nil
}
