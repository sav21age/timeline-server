package email

import "errors"

var MsgMailServiceUnavailable = "mail service unavailable"

var ErrMailServiceUnavailable = errors.New(MsgMailServiceUnavailable)

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}