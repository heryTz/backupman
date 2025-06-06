//go:build test_integration

package mail_test

import (
	"testing"

	"github.com/herytz/backupman/core/notifier/mail"
	"github.com/stretchr/testify/assert"
)

func TestSendMailSuccess(t *testing.T) {
	mailer := mail.NewStdMailNotifier("localhost", 1026, "", "", "")
	err := mailer.Send(mail.MailNotifyInput{
		Recipients: []mail.Recipient{
			{Name: "John Doe", Email: "johndoe3@yopmail.fr"},
		},
		Subject: "Test Subject",
		Message: "Test Message",
	})
	assert.NoError(t, err)
}
