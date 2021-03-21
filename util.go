package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Println(result.Responses[i].FullTextAnnotation.Text)
		strList = append(strList, result.Responses[i].FullTextAnnotation.Text)
	}
	return strList
}
