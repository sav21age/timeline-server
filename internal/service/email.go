package service

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"text/template"
	"timeline/config"
	"timeline/internal/domain"
	"timeline/pkg/email"
)

const htmlTplEmailDir = "../../template/email/"

type EmailService struct {
	sender *email.EmailSender
	cfg    *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	s := email.NewEmailSender(cfg.From, cfg.Email.Password, cfg.Email.SMTP.Host, cfg.Email.SMTP.Port)
	
	return &EmailService{
		sender: s,
		cfg:    cfg,
	}
}

//go:generate mockgen -source=email.go -destination=mock/email.go

type EmailInterface interface {
	AccountActivateEmail(domain.User) error
	PasswordRecoveryEmail(domain.User) error
}

func generateBody(htmlTplPath string, data interface{}) (string, error) {
	var b bytes.Buffer

	_, mainFilePath, _, _ := runtime.Caller(1)
	mainDirPath := filepath.Dir(mainFilePath)
	htmlTplFullPath := filepath.Join(mainDirPath, htmlTplPath)

	t, err := template.ParseFiles(htmlTplFullPath)
	if err != nil {
		// return "", errors.New("error parse template file")
		return "", err
	}

	t.Execute(&b, data)

	return b.String(), err
}

//--

func (s *EmailService) AccountActivateEmail(user domain.User) error {
	activationLink := fmt.Sprintf("%s/auth/account/activate/%s", s.cfg.Client.Url, user.Activation.Code)
	htmlTplPath := htmlTplEmailDir + "activate_account.html"
	subject := "Please, confirm email address"

	data := struct {
		Name string
		Link string
	}{
		Name: user.Username,
		Link: activationLink,
	}

	body, err := generateBody(htmlTplPath, data)
	if err != nil {
		return err
	}

	// return email.SendMailer(msg, s.cfg)
	return s.sender.Send(
		email.EmailMessage{
			To:      user.Email,
			Subject: subject,
			Body:    body,
		})
}

//--

func (s *EmailService) PasswordRecoveryEmail(user domain.User) error {
	ResetLink := fmt.Sprintf("%s/auth/password/reset/%s", s.cfg.Client.Url, user.Recovery.Code)
	htmlTplPath := htmlTplEmailDir + "recovery_password.html"
	subject := "Password recovery"

	data := struct {
		Name string
		Link string
	}{
		Name: user.Username,
		Link: ResetLink,
	}

	body, err := generateBody(htmlTplPath, data)
	if err != nil {
		return err
	}

	// return email.SendMailer(msg, s.cfg)
	return s.sender.Send(
		email.EmailMessage{
			To:      user.Email,
			Subject: subject,
			Body:    body,
		})
}
