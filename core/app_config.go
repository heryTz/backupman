package core

type GeneralConfig struct {
	AppUrl     string
	BackupCron string
}

type DriveConfig interface{}
type LocalDriveConfig struct {
	Label  string
	Folder string
}
type GoogleDriveConfig struct {
	Label          string
	Folder         string
	ServiceAccount string
}

type DataSourceConfig interface{}
type MysqlDataSourceConfig struct {
	Label     string
	TmpFolder string
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Tls       string
}

type DbConfig interface{}
type MysqlDbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Tls      string
}
type MemoryDbConfig struct{}

type MailNotifierDestinationConfig struct {
	Name  string
	Email string
}
type MailNotifierConfig struct {
	Enabled      bool
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
	SmtpCrypto   string
	Destinations []MailNotifierDestinationConfig
}
type NotifierConfig struct {
	Mail MailNotifierConfig
}

type RetentionConfig struct {
	Enabled bool
	Days    int
}

type WebhookConfig struct {
	Name  string
	Url   string
	Token string
}

type AppConfig struct {
	General     GeneralConfig
	ApiKeys     []string
	Drives      []DriveConfig
	DataSources []DataSourceConfig
	Db          DbConfig
	Notifiers   NotifierConfig
	Retention   RetentionConfig
	Webhooks    []WebhookConfig
}
