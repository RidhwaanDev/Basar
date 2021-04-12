package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ocr result
type OCRResult struct {
	Responses []Response `json:"responses"`
}

type Response struct {
	FullTextAnnotation Annotation `json:"fullTextAnnotation"`
}

type Annotation struct {
	Text string `json:"text"`
}

func printFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := ioutil.ReadAll(file)
	fmt.Print(string(b))
}

func ParseJSONFile(fileName string) []string {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result OCRResult

	json.Unmarshal(byteValue, &result)

	var strList []string
	for i := 0; i < len(result.Responses); i++ {
		// fmt.Println(result.Responses[i].FullTextAnnotation.Text)
		strList = append(strList, result.Responses[i].FullTextAnnotation.Text)
	}
	return strList
}

func GenRandomID() string {
	return uuid.NewString()
}

func CleanOutFinalResult(fileName string) {}

func CleanDownloadedFiles(prefix string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) && filepath.Ext(f.Name()) != ".txt" {
			err := os.Remove(f.Name())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
