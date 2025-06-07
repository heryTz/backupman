package drive

type DriveFile struct {
	Path string
}

type Drive interface {
	Upload(srcPath string) (DriveFile, error)
	Delete(srcPath string) error
	GetLabel() string
	GetProvider() string
	Health() error
}
