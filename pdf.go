package main

import (
	"github.com/signintech/gopdf"
)

func main() {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	pdf.Cell(nil, "您好")
	pdf.WritePdf("hello.pdf")

}
