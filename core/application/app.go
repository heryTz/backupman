package application

import (
	"log"

	"github.com/herytz/backupman/core/dao"
	"github.com/herytz/backupman/core/dao/memory"
	"github.com/herytz/backupman/core/dao/mysql"
	"github.com/herytz/backupman/core/dao/postgres"
	"github.com/herytz/backupman/core/dao/sqlite"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/dumper"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/mailer"
	"github.com/herytz/backupman/core/notifier"
)

const APP_MODE_CLI = "cli"
const APP_MODE_WEB = "web"

type Webhook struct {
	Name  string
	Url   string
	Token string
}

type App struct {
	Version struct {
		Version   string
		CommitSHA string
		BuildDate string
	}
	Mode string
	Http struct {
		AppUrl    string
		ApiKeys   []string
		BackupJob struct {
			Enabled bool
			Cron    string
		}
	}
	Drives    []drive.Drive
	Dumpers   []dumper.Dumper
	Db        dao.Dao
	Notifiers []notifier.Notifier
	Retention struct {
		Enabled bool
		Days    int
	}
}

func NewApp(config AppConfig) *App {
	app := App{}

	app.Version.Version = config.Version.Version
	app.Version.CommitSHA = config.Version.CommitSHA
	app.Version.BuildDate = config.Version.BuildDate

	drives := make([]drive.Drive, len(config.Drives))
	for i, driveConfig := range config.Drives {
		switch config := driveConfig.(type) {
		case LocalDriveConfig:
			drives[i] = drive.NewLocalDrive(config.Label, config.Folder)
		case GoogleDriveConfig:
			drives[i] = drive.NewGoogleDrive(config.Label, config.Folder, config.ClientSecretFile, config.TokenFile)
		default:
			log.Fatal("Unsupported drive type")
		}
	}

	dumpers := make([]dumper.Dumper, len(config.DataSources))
	for i, dataSourceConfig := range config.DataSources {
		switch config := dataSourceConfig.(type) {
		case MysqlDataSourceConfig:
			dumpers[i] = dumper.NewMysqlDumper(
				config.Label,
				config.TmpFolder,
				config.Host,
				config.Port,
				config.User,
				config.Password,
				config.Database,
				config.Tls,
			)
		case PostgresDataSourceConfig:
			dumpers[i] = dumper.NewPostgresDumper(
				config.Label,
				config.TmpFolder,
				config.Host,
				config.Port,
				config.User,
				config.Password,
				config.Database,
				config.Tls,
			)
		default:
			log.Fatal("Unsupported database type")
		}
	}

	db := dao.Dao{}
	switch config := config.Db.(type) {
	case MysqlDbConfig:
		dbConn, err := lib.NewMysqlConnection(
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.Database,
			config.Tls,
		)
		if err != nil {
			log.Fatal(err)
		}
		db.Backup = mysql.NewBackupDaoMysql(dbConn)
		db.DriveFile = mysql.NewDriveFileDaoMysql(dbConn)
		db.Health = lib.NewHealthMysql(dbConn)
	case PostgresDbConfig:
		dbConn, err := lib.NewPostgresConnection(
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.Database,
			config.Tls,
		)
		if err != nil {
			log.Fatal(err)
		}
		db.Backup = postgres.NewBackupDaoPostgres(dbConn)
		db.DriveFile = postgres.NewDriveFileDaoPostgres(dbConn)
		db.Health = lib.NewHealthPostgres(dbConn)
	case SqliteDbConfig:
		dbConn, err := lib.NewSqliteConnection(config.DbPath)
		if err != nil {
			log.Fatal(err)
		}
		db.Backup = sqlite.NewBackupDaoSqlite(dbConn)
		db.DriveFile = sqlite.NewDriveFileDaoSqlite(dbConn)
		db.Health = lib.NewHealthSqlite(dbConn)
	case MemoryDbConfig:
		memoryDb := memory.NewMemoryDb()
		db.Backup = memory.NewBackupDaoMemory(memoryDb)
		db.DriveFile = memory.NewDriveFileDaoMemory(memoryDb)
		db.Health = lib.MockUpHelthChecker{}
	default:
		log.Fatal("Unsupported dao type")
	}

	notifiers := []notifier.Notifier{}

	mailConfig := config.Notifiers.Mail
	if mailConfig.Enabled {
		mailerTransport := mailer.NewStdMailer(
			mailConfig.SmtpHost,
			mailConfig.SmtpPort,
			mailConfig.SmtpUser,
			mailConfig.SmtpPassword,
			mailConfig.SmtpCrypto,
		)
		mailDestinations := []mailer.Recipient{}
		for _, destination := range mailConfig.Destinations {
			mailDestinations = append(mailDestinations, mailer.Recipient{
				Name:  destination.Name,
				Email: destination.Email,
			})
		}
		notifiers = append(notifiers, notifier.NewMailNotifier(mailerTransport, db, mailDestinations))
	}

	webhookConfig := config.Notifiers.Webhooks
	if len(webhookConfig) > 0 {
		webhookNotifierConfigs := []notifier.WebhookNotifierConfig{}
		for _, config := range webhookConfig {
			webhookNotifierConfigs = append(webhookNotifierConfigs, notifier.WebhookNotifierConfig{
				Name:  config.Name,
				Url:   config.Url,
				Token: config.Token,
			})
		}
		notifiers = append(notifiers, notifier.NewWebhookNotifier(webhookNotifierConfigs, db))
	}

	app.Dumpers = dumpers
	app.Drives = drives
	app.Db = db
	app.Notifiers = notifiers

	app.Http.ApiKeys = config.Http.ApiKeys
	app.Http.AppUrl = config.Http.AppUrl
	app.Http.BackupJob.Enabled = config.Http.BackupJob.Enabled
	app.Http.BackupJob.Cron = config.Http.BackupJob.Cron

	app.Retention.Enabled = config.Retention.Enabled
	app.Retention.Days = config.Retention.Days

	return &app
}
