package mail

type MockMailNotifier struct{}

func NewMockMailNotifier() *MockMailNotifier {
	return &MockMailNotifier{}
}

func (m *MockMailNotifier) Send(input MailNotifyInput) error {
	return nil
}

func (m *MockMailNotifier) Health() error {
	return nil
}
