package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"log"
)

func main() {
	adapter := NewSmtpAdapter("smtp.mail.ru",
		"ilyayaldinov@inbox.ru",
		"RvSXC4wkVdunpmqkQE89",
		"ilyayaldinov@inbox.ru", 587)

	if err := adapter.Send("ilyayaldinov@inbox.ru", "ку", "ку"); err != nil {
		log.Fatal(err)
	}
}

type SmtpAdapter struct {
	dialer *gomail.Dialer
	email  string
}

func NewSmtpAdapter(host, login, password, email string, port int) *SmtpAdapter {
	return &SmtpAdapter{
		email:  email,
		dialer: gomail.NewDialer(host, port, login, password),
	}
}

func (smtpAdapter *SmtpAdapter) Send(target string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", smtpAdapter.email)
	message.SetHeader("To", target)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)
	return smtpAdapter.dialer.DialAndSend(message)
}

func (smtpAdapter *SmtpAdapter) MultipleSend(subject string, body string, targets ...string) error {
	errList := make([]error, 0, 2)
	for _, target := range targets {
		if err := smtpAdapter.Send(target, subject, body); err != nil {
			errList = append(errList, err)
		}
	}
	return smtpAdapter.processErrors(errList)
}

func (smtpAdapter *SmtpAdapter) processErrors(errs []error) error {
	errorMessage := "Ошибка отправки писем: "
	if len(errs) == 0 {
		return nil
	}
	for _, error := range errs {
		errorMessage = fmt.Sprint(errorMessage, error.Error(), " ")
	}

	return errors.New(errorMessage)

}
