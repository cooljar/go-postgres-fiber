package utils

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"net"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const forceDisconnectAfter = time.Second * 5

type Mailer struct {
	Subject string
	Destination string
	ViewPath string
	Data map[string]interface{}
}

func NewMailer(subject, dest, viewPath string, data map[string]interface{}) *Mailer {
	return &Mailer{
		Subject: subject,
		Destination: dest,
		ViewPath: viewPath,
		Data: data,
	}
}

// IsEmailAddressValid validate mail host.
func (m *Mailer) IsEmailAddressValid() (valid bool, err error) {
	_, host := split(m.Destination)
	mx, err := net.LookupMX(host)
	if err != nil {
		return false, err
	}
	client, err := dialTimeout(fmt.Sprintf("%s:%d", mx[0].Host, 25), forceDisconnectAfter)
	if err != nil {
		return false, err
	}
	client.Close()

	return true, nil
}

func (m *Mailer) SendMail() error {
	var pathLayout = path.Join("assets", "html", "layout", "mail.html")
	var pathContent = m.ViewPath

	var tpl bytes.Buffer

	templateEmail, err := template.ParseFiles(pathLayout, pathContent)
	if err != nil {
		return err
	}

	err = templateEmail.Execute(&tpl, m.Data)
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SMTP_SENDER_NAME"))
	mailer.SetHeader("To", m.Destination)
	//mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Email Validation Secret Code At - Cooljar Apps")
	mailer.SetBody("text/html", tpl.String())
	//mailer.Attach("./sample.png")

	portInt, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		portInt,
		os.Getenv("SMTP_AUTH_EMAIL"),
		os.Getenv("SMTP_AUTH_PASSWORD"),
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

// DialTimeout returns a new Client connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".
func dialTimeout(addr string, timeout time.Duration) (*smtp.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	t := time.AfterFunc(timeout, func() { conn.Close() })
	defer t.Stop()

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}


func split(email string) (account, host string) {
	i := strings.LastIndexByte(email, '@')
	// If no @ present, not a valid email.
	if i < 0 {
		return
	}
	account = email[:i]
	host = email[i+1:]
	return
}
