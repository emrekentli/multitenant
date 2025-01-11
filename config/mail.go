package config

import (
	"github.com/valyala/bytebufferpool"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Mail struct {
	*mail.SMTPServer
	*mail.SMTPClient
	Host        string
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	View        *ViewConfig
	Port        int
}

func (m *Mail) Send(to string, subject string, body string, cc string, from string) error {
	if m.SMTPServer == nil {
		m.SetupMailer()
	}
	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom(from).AddTo(to).SetSubject(subject) // nolint:wsl
	if cc != "" {                                     // nolint:wsl
		email.AddCc(cc)
	}
	email.SetBody(mail.TextHTML, body) // nolint:wsl

	// Call Send and pass the client
	err := email.Send(m.SMTPClient)

	if err != nil {
		return err
	} else {
		log.Println("Email Sent")
	}
	return nil
}

func (m *Mail) PrepareHtml(view string, body fiber.Map) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	// app.Settings.Views.Render
	if err := m.View.Template.TemplateEngine.Render(buf, view, body); err != nil {
		// handle err
	}
	return buf.String()
}

func (m *Mail) SetupMailer() {
	m.Host = "smtp.mailtrap.io"
	m.Username = "821c8fc0bb1e19"
	m.Password = "24edfcaf91afbc"
	m.Encryption = "tls"
	m.FromAddress = "example@gmail.com"
	m.FromName = "mail name"
	m.Port = 2525

	var err error
	m.SMTPServer = mail.NewSMTPClient()
	m.SMTPServer.Host = m.Host
	m.SMTPServer.Port = m.Port
	m.SMTPServer.Username = m.Username
	m.SMTPServer.Password = m.Password
	if m.Encryption == "tls" {
		m.SMTPServer.Encryption = mail.EncryptionTLS
	} else {
		m.SMTPServer.Encryption = mail.EncryptionSSL
	}

	// Variable to keep alive connection
	m.SMTPServer.KeepAlive = false

	// Timeout for connect to SMTP Server
	m.SMTPServer.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	m.SMTPServer.SendTimeout = 10 * time.Second
	m.SMTPClient, err = m.SMTPServer.Connect()
	if err != nil {
		log.Print(err)
	}
}
