package main

import (
	"fmt"
)

// parse the OCR results. Remove non-Arabic characters, and symbols
func ParseText(s string) string {

}

func ApplyAllah(s string) string {
	if s == "ال له" {
		return "الله"
	}
}

func ApplyRasul(s string) string {
	if s == "رشل" {
		return "رسل"
	}
}
