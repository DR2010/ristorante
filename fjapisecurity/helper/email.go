package helper

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Mail is for email
type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

// SMTPServer is for SmtpServer
type SMTPServer struct {
	host string
	port string
}

// ServerName is ServerName
func (s *SMTPServer) ServerName() string {
	return s.host + ":" + s.port
}

// BuildMessage is
func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

// SendEmail test
func SendEmail(destination string, message string) {

	mail := Mail{}
	mail.senderId = "festajuninacanberra@gmail.com"
	mail.toIds = []string{destination}
	mail.subject = "Machadodaniel.com password reset"
	mail.body = message

	messageBody := mail.BuildMessage()

	// port 465 SSL and 587 TLS
	log.Println("SMTPServer := SMTPServer{host: smtp.gmail.com, port: 465}")

	SMTPServer := SMTPServer{host: "smtp.gmail.com", port: "465"}
	// SMTPServer := SMTPServer{host: "smtp.gmail.com", port: "587"}

	log.Println(SMTPServer.host)
	//build an auth
	log.Println("auth := smtp.PlainAuth(, mail.senderId, xegbyiwwkiijoysm, SMTPServer.host)")
	auth := smtp.PlainAuth("", mail.senderId, "xegbyiwwkiijoysm", SMTPServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPServer.host,
	}

	conn, err := tls.Dial("tcp", SMTPServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, SMTPServer.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
