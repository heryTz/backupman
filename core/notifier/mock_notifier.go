package notifier

type MockNotifier struct{}

func (m *MockNotifier) BackupReport(backupId string) error {
	return nil
}

func (m *MockNotifier) Health() error {
	return nil
}

func (m *MockNotifier) GetName() string {
	return "mock"
}
