package drive

type DriveFile struct {
	Path string
}

type Drive interface {
	Upload(srcPath string) (DriveFile, error)
	GetLabel() string
	GetProvider() string
}
