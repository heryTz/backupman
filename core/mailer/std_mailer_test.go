//go:build test_integration

package mailer_test

import (
	"testing"

	"github.com/herytz/backupman/core/mailer"
	"github.com/stretchr/testify/assert"
)

func TestSendMailSuccess(t *testing.T) {
	mailerTransport := mailer.NewStdMailer("localhost", 1026, "", "", "")
	err := mailerTransport.Send(mailer.MailerInput{
		Recipients: []mailer.Recipient{
			{Name: "John Doe", Email: "johndoe3@yopmail.fr"},
		},
		Subject: "Test Subject",
		Message: "Test Message",
	})
	assert.NoError(t, err)
}
