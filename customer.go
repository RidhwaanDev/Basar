package main

import (
	"fmt"
	"strings"
)

type customer struct {
	name     string
	contact  string
	madrasah string
}

func (c customer) hasEmail() bool {
	return strings.Contains(c.contact, "@")
}
