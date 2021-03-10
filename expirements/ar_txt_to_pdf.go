package main

// TODO make this library with docker

// headless chrome lol
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ConvertTextToPDF(s string) {
	// chrome --headless --disable-gpu --print-to-pdf text.txt

	// convert the string to bytes
	b := []byte(s)

	err := ioutil.WriteFile("text.txt", b, 0644)
	check(err)

	prg := "chrome"
	arg1 := "--headless"
	arg2 := "--disable-gpu"
	arg3 := "--print-to-pdf"
	val1 := "text.txt"
	cmd := exec.Command(prg, arg1, arg2, arg3, val1)

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("an error in pdf convert" + "   " + stderr.String())
		log.Fatal(err)
	} else {
		fmt.Println("successfully converted .txt to .pdf")
	}
}
