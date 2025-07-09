package notifier

import (
	"github.com/herytz/backupman/core/dao"
	"github.com/herytz/backupman/core/mailer"
	"github.com/herytz/backupman/core/notifier/message"
)

type MailNotifier struct {
	Mailer     mailer.Mailer
	Db         dao.Dao
	Recipients []mailer.Recipient
}

func NewMailNotifier(mailer mailer.Mailer, db dao.Dao, recipients []mailer.Recipient) *MailNotifier {
	return &MailNotifier{Mailer: mailer, Db: db, Recipients: recipients}
}

func (m *MailNotifier) BackupReport(backupId string) error {
	backup, err := m.Db.Backup.ReadFullById(backupId)
	if err != nil {
		return err
	}

	msg, err := message.BackupReportMail(backup)
	if err != nil {
		return err
	}

	var input mailer.MailerInput
	input.Recipients = m.Recipients
	input.Subject = "Backup Report"
	input.Message = msg

	return m.Mailer.Send(input)
}

func (m *MailNotifier) Health() error {
	return m.Mailer.Health()
}

func (m *MailNotifier) GetName() string {
	return "Mail"
}
