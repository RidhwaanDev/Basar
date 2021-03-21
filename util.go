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
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var result OCRResult

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &result)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	var strList []string
	for i := 0; i < len(result.Responses); i++ {
		fmt.Println("OCR Text:" + result.Responses[i].FullTextAnnotation.Text)
		strList = append(strList, result.Responses[i].FullTextAnnotation.Text)
	}
	return strList
}
