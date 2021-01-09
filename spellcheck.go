package main

import (
	"fmt"
)

// parse the OCR results. Remove non-Arabic characters, and symbols
func ParseText(s string) string {

}

func ApplyRule(s string) string {
	switch s {
	case "ال له":
		return "الله" 
	case "رشل":
		return "رسل"
        }
}

