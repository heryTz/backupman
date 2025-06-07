package mail

type Recipient struct {
	Name  string
	Email string
}

type MailNotifyInput struct {
	Recipients []Recipient
	Subject    string
	Message    string
}

type MailNotifier interface {
	Send(input MailNotifyInput) error
	Health() error
}
