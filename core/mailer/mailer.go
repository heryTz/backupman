package mailer

type Recipient struct {
	Name  string
	Email string
}

type MailerInput struct {
	Recipients []Recipient
	Subject    string
	Message    string
}

type Mailer interface {
	Send(input MailerInput) error
	Health() error
}
