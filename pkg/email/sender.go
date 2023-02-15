package email

import (
	"errors"
	"strconv"
	"syscall"

	"github.com/go-gomail/gomail"
	"github.com/rs/zerolog/log"
)

type EmailSender struct {
	from     string
	password string
	host     string
	port     string
}

func NewEmailSender(from, password, host, port string) *EmailSender {
	return &EmailSender{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (s *EmailSender) Send(msg EmailMessage) error {
	port, err := strconv.Atoi(s.port)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", msg.Body)

	dialer := gomail.NewDialer(s.host, port, s.from, s.password)
	err = dialer.DialAndSend(m)
	if errors.Is(err, error(syscall.Errno(10061))) {
		log.Warn().Msg(MsgMailServiceUnavailable)
		return ErrMailServiceUnavailable
	}

	return err
}
