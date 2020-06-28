package email

import (
	"errors"
	"log"
	"net/smtp"
	"os"
)

// ErrSendMail is returned when the server could not generate or send
// a mail
var ErrSendMail = errors.New("Could not send mail")

var (
	smtpServerEnv string

	// BotEmailEnv is an environment variable used to set server's bot email
	// address, bot should be a no-reply email address used for email confirmation,
	// password reset etc.
	BotEmailEnv string
)

// Send is used to send a mail to a user
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
