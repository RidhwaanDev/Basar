package main

import (
	"github.com/signintech/gopdf"
)

type pdf struct {
	location string // location on disk
	name     string
	pages    int
	size     int64
}

// number of pages in the PDF
func (pdf p) GetPageCount(ch <-chan int) {}

// mb size of PDF, usually when client uploads pdf, we will already know that
func (pdf p) GetSize(ch <-chan int64) {}
