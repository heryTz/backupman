package notifier

type Notifier interface {
	BackupReport(backupId string) error
	Health() error
	GetName() string
}
