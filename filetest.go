package main

import (
	"fmt"
	"path"
	"strings"
)

func main() {
	str := "test.jpg"
	fmt.Println(strings.TrimSuffix(str, path.Ext(str)))
}
