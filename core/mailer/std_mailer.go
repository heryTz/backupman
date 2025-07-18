package mailer

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// - SSL is the predecessor of TLS
// - We can etablish a secure connection using TLS or SSL
// -- At the beginning of the connection, using SMTPS over TLS or SMTPS over SSL (port 465)
// -- After the connection is established, using STARTTLS to upgrade the connection to TLS (port 587)

type StdMailer struct {
	smtpHost   string
	smtpPort   int
	smtpUser   string
	smtpPass   string
	smtpCrypto string
}

func NewStdMailer(smtpHost string, smtpPort int, smtpUser, smtpPass, smtpCrypto string) *StdMailer {
	return &StdMailer{
		smtpHost:   smtpHost,
		smtpPort:   smtpPort,
		smtpUser:   smtpUser,
		smtpPass:   smtpPass,
		smtpCrypto: smtpCrypto,
	}
}

func (m *StdMailer) Send(input MailerInput) error {
	var err error
	var client *smtp.Client

	switch strings.ToLower(m.smtpCrypto) {

	case "tls":
		tlsConfig := &tls.Config{
			ServerName: m.smtpHost,
		}
		client, err = smtp.DialStartTLS(fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort), tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP tls server => %v", err)
		}
		defer client.Close()

		auth := sasl.NewLoginClient(m.smtpUser, m.smtpPass)
		err = client.Auth(auth)
		if err != nil {
			return fmt.Errorf("failed to authenticate => %v", err)
		}

	case "ssl":
		tlsConfig := &tls.Config{
			ServerName:         m.smtpHost,
			InsecureSkipVerify: true,
		}
		client, err = smtp.DialTLS(fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort), tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP ssl server => %v", err)
		}
		defer client.Close()

		auth := sasl.NewLoginClient(m.smtpUser, m.smtpPass)
		err = client.Auth(auth)
		if err != nil {
			return fmt.Errorf("failed to authenticate => %v", err)
		}

	case "":
		client, err = smtp.Dial(fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort))
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server => %v", err)
		}
		defer client.Close()

	default:
		return fmt.Errorf("unsupported crypto type: %s", m.smtpCrypto)
	}

	var to []string
	for _, recepient := range input.Recipients {
		to = append(to, recepient.Email)
	}

	from := "no-replay@fake.com"
	if m.smtpUser != "" {
		from = m.smtpUser
	}

	msg := strings.NewReader(
		"From: " + from + "\r\n" +
			"To: " + strings.Join(to, ", ") + "\r\n" +
			"Subject: " + input.Subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			input.Message + "\r\n",
	)

	err = client.SendMail(from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email => %v", err)
	}
	return nil
}

func (m *StdMailer) Health() error {
	return m.Send(MailerInput{
		Subject: "Health Check",
		Message: "This is a health check email from the StdMailer.",
		Recipients: []Recipient{
			{Name: "Health Check", Email: "healh.check@yopmail.fr"},
		},
	})
}
