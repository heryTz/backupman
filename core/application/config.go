package application

type HttpConfig struct {
	AppUrl    string
	ApiKeys   []string
	BackupJob struct {
		Enabled bool
		Cron    string
	}
}

type DriveConfig interface{}
type LocalDriveConfig struct {
	Label  string
	Folder string
}
type GoogleDriveConfig struct {
	Label            string
	Folder           string
	ClientSecretFile string
	TokenFile        string
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
type PostgresDataSourceConfig struct {
	Label     string
	TmpFolder string
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Tls       bool
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
type PostgresDbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Tls      bool
}
type SqliteDbConfig struct {
	DbPath string
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

type WebhookNotifierConfig struct {
	Name  string
	Url   string
	Token string
}

type NotifierConfig struct {
	Mail     MailNotifierConfig
	Webhooks []WebhookNotifierConfig
}

type RetentionConfig struct {
	Enabled bool
	Days    int
}

type VersionConfig struct {
	Version   string
	CommitSHA string
	BuildDate string
}

type AppConfig struct {
	Http        HttpConfig
	Drives      []DriveConfig
	DataSources []DataSourceConfig
	Db          DbConfig
	Notifiers   NotifierConfig
	Retention   RetentionConfig
	Version     VersionConfig
}
