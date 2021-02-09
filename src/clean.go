package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

func CleanUpUploadsFolder() {
	items, err := ioutil.ReadDir("uploads")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		// fmt.Println(item.Name())
		err := os.Remove("uploads/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}
	// remove the pdf file
	os.Remove("pdf_to_convert.pdf")
}

// remove everything from results dir
func CleanUpResultsFolder() {
	items, err := ioutil.ReadDir("results")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		// fmt.Println(item.Name())
		err := os.Remove("results/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func removeEmptyLines() {
	file, err := os.Open("final_output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// final_output.txt
func createFinalOutputTextFile() {
	items, err := ioutil.ReadDir("results")

	names := make([]string, 1)
	if err != nil {
		fmt.Println(err)
	}
	var itemCount int
	for _, item := range items {
		itemCount++
		trimmed := strings.Trim(item.Name(), path.Ext(item.Name()))
		trimmed2 := strings.Trim(trimmed, "output_file-")
		fmt.Println(trimmed2)
		names = append(names, trimmed2)
		err := os.Rename("results/"+item.Name(), "results/"+trimmed2+".txt")
		// 	err := os.Remove("results/" + item.Name())
		if err != nil {
			fmt.Println(err)
		}
	}

	sort.Slice(names, func(i, j int) bool {
		numA, _ := strconv.Atoi(names[i])
		numB, _ := strconv.Atoi(names[j])
		return numA < numB
	})

	fmt.Printf("After: %v\n", names)
	finalTextFile, err := os.OpenFile("final_output.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer finalTextFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	// read all the data from the output files and put them into final_output.txt

	for i := 0; i <= itemCount; i++ {
		outFile, err := os.Open("results/" + names[i] + ".txt")

		if err != nil {
			fmt.Println("error in opening the results file")
		}

		bytes, err := ioutil.ReadAll(outFile)
		if err != nil {
			fmt.Println("error in reading the bytes from output file")
		}

		n, err := io.WriteString(finalTextFile, string(bytes))

		if err != nil {
			fmt.Println("error in writing the bytes to the output file")
		} else {
			// fmt.Printf("%d bytes written to %s", n, names[i]+".txt")
		}
		outFile.Close()
	}

	removeEmptyLines()

}
