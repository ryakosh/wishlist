package email

import (
	"log"
	"net/smtp"
	"os"
)

var (
	smtpServerEnv string
	BotEmailEnv   string
)

func Send(email string, to string, sub string, msg string) error {
	tos := []string{to}
	body := []byte("To: " + tos[0] + "\r\n" +
		"Subject: " + sub + "\r\n" +
		"Content-Type: text/html; charset=UTF-8" + "\r\n" +
		"\r\n" +
		msg,
	)

	return smtp.SendMail(smtpServerEnv, nil, email, tos, body)
}

func init() {
	smtpServerEnv = os.Getenv("WISHLIST_SMTPSERVER")
	if len(smtpServerEnv) == 0 {
		log.Fatal("error: 'WISHLIST_SMTPSERVER' must be set")
	}

	BotEmailEnv = os.Getenv("WISHLIST_BOTEMAIL")
	if len(BotEmailEnv) == 0 {
		log.Fatal("error: 'WISHLIST_BOTEMAIL' must be set")
	}
}
