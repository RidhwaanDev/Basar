package main

import (
	"fmt"
	"sort"
)

func main() {
	str := []string{
		"output_file-0.txt",
		"output_file-1.txt",
		"output_file-10.txt",
		"output_file-11.txt",
		"output_file-12.txt",
		"output_file-13.txt",
		"output_file-14.txt",
		"output_file-15.txt",
		"output_file-16.txt",
		"output_file-17.txt",
		"output_file-18.txt",
		"output_file-2.txt",
		"output_file-3.txt",
		"output_file-4.txt",
		"output_file-5.txt",
		"output_file-6.txt",
		"output_file-7.txt",
		"output_file-8.txt",
		"output_file-9.txt",
	}

	sort.Strings(str)
	fmt.Println(str)

}
