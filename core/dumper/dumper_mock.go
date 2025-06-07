package dumper

type DumperMock struct{}

func (d *DumperMock) Dump() (string, error) {
	return "./dumper_mock_db", nil
}

func (d *DumperMock) GetLabel() string {
	return "dumper_mock"
}

func (d *DumperMock) Health() error {
	return nil
}
