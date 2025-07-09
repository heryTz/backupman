package mailer

type MockMailNotifier struct{}

func NewMockMailNotifier() *MockMailNotifier {
	return &MockMailNotifier{}
}

func (m *MockMailNotifier) Send(input MailerInput) error {
	return nil
}

func (m *MockMailNotifier) Health() error {
	return nil
}
