package main

import (
	"fmt"
	"net/smtp"
)

// pop3 only downloads the contents of your inbox folder, does not sync your folders or any other ermails, also deletes the email from the server
// imap syncs, plus caches it on the server, plus syncs between your devices
// smtp is a protcol. simple mail transfer protocl (sending mail to people)
// smtp uses TCP protocl underneath

func main() {
	from := "ridhwaan.any@gmail.com"
	password := "mypass"
	to := []string{
		"noorsyed1a@gmail.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte("This is a test email message.")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("succesfully")
}
