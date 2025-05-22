package drive

type DriveMock struct{}

func (d *DriveMock) Upload(srcPath string) (DriveFile, error) {
	return DriveFile{
		Path: "./drive_mock/file.txt",
	}, nil
}

func (d *DriveMock) Delete(dstPath string) error {
	return nil
}

func (d *DriveMock) GetLabel() string {
	return "drive_mock"
}

func (d *DriveMock) GetProvider() string {
	return "mock"
}
