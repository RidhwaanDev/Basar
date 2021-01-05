package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type pdf struct {
	location string // location on disk
	name     string
	pages    int
	size     int64
}

// number of pages in the PDF
func (p pdf) GetPageCount(ch <-chan int) {}

// mb size of PDF, usually when client uploads pdf, we will already know that
func (p pdf) GetSize(ch <-chan int64) {}

const (
	filename = "input.txt"
)

//func main() {
//	file, err := os.Create("input.txt")
//	defer os.Remove("input.txt")
//
//	data := "هيعشف هيشغج شهرغج فذلجة هاسلثعف ةىچآشد غ لثقشهدنبدصزخ،مصذ نغ"
//	check(err)
//	err = ioutil.WriteFile("input.txt", []byte(data), 0666)
//	ConvertTextToPDF(file)
//}

// headless chrome lol
func ConvertTextToPDF(file *os.File) *os.File {
	// chrome --headless --disable-gpu --print-to-pdf input.txt
	info, _ := file.Stat()
	defer os.Remove(info.Name())

	chromeExec, err := exec.LookPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome")
	check(err)

	cmd := &exec.Cmd{
		Path:   chromeExec,
		Args:   []string{chromeExec, "--headless", "--disable-gpu", "--print-to-pdf", "input.txt"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err = cmd.Run(); err != nil {
		fmt.Println("an error in pdf convert")
		log.Fatal(err)
		return nil
	}

	var outFile *os.File

	fmt.Println("successfully converted .txt to .pdf")

	if outFile, err = os.Open("output.pdf"); err != nil {
		check(err)
	}

	return outFile
}
