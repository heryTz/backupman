package application

import (
	"log"

	"github.com/herytz/backupman/core/dao"
	"github.com/herytz/backupman/core/dao/memory"
	"github.com/herytz/backupman/core/dao/mysql"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/dumper"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/notifier/mail"
)

const APP_MODE_CLI = "cli"
const APP_MODE_WEB = "web"

type Webhook struct {
	Name  string
	Url   string
	Token string
}

type App struct {
	Mode string
	Http struct {
		AppUrl    string
		ApiKeys   []string
		BackupJob struct {
			Enabled bool
			Cron    string
		}
	}
	Drives       []drive.Drive
	Dumpers      []dumper.Dumper
	Db           dao.Dao
	Notification struct {
		Mail struct {
			Enabled      bool
			Destinations []mail.Recipient
		}
	}
	Notifiers struct {
		Mail mail.MailNotifier
	}
	Retention struct {
		Enabled bool
		Days    int
	}
	Webhooks []Webhook
}

func NewApp(config AppConfig) *App {
	app := App{}

	drives := make([]drive.Drive, len(config.Drives))
	for i, driveConfig := range config.Drives {
		switch config := driveConfig.(type) {
		case LocalDriveConfig:
			drives[i] = drive.NewLocalDrive(config.Label, config.Folder)
		case GoogleDriveConfig:
			drives[i] = drive.NewGoogleDrive(config.Label, config.Folder, config.ServiceAccount)
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
		default:
			log.Fatal("Unsupported database type")
		}
	}

	db := dao.Dao{}
	switch config := config.Db.(type) {
	case MysqlDbConfig:
		dbConn, err := lib.NewConnection(
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
	case MemoryDbConfig:
		memoryDb := memory.NewMemoryDb()
		db.Backup = memory.NewBackupDaoMemory(memoryDb)
		db.DriveFile = memory.NewDriveFileDaoMemory(memoryDb)
	default:
		log.Fatal("Unsupported dao type")
	}

	mailConfig := config.Notifiers.Mail
	if mailConfig.Enabled {
		app.Notifiers.Mail = mail.NewStdMailNotifier(
			mailConfig.SmtpHost,
			mailConfig.SmtpPort,
			mailConfig.SmtpUser,
			mailConfig.SmtpPassword,
			mailConfig.SmtpCrypto,
		)
		mailDestinations := []mail.Recipient{}
		for _, destination := range mailConfig.Destinations {
			mailDestinations = append(mailDestinations, mail.Recipient{
				Name:  destination.Name,
				Email: destination.Email,
			})
		}

		app.Notification.Mail.Destinations = mailDestinations
		app.Notification.Mail.Enabled = mailConfig.Enabled
	}

	app.Webhooks = []Webhook{}
	for _, wh := range config.Webhooks {
		app.Webhooks = append(app.Webhooks, Webhook{
			Name:  wh.Name,
			Url:   wh.Url,
			Token: wh.Token,
		})
	}

	app.Dumpers = dumpers
	app.Drives = drives
	app.Db = db

	app.Http.ApiKeys = config.Http.ApiKeys
	app.Http.AppUrl = config.Http.AppUrl
	app.Http.BackupJob.Enabled = config.Http.BackupJob.Enabled
	app.Http.BackupJob.Cron = config.Http.BackupJob.Cron

	app.Retention.Enabled = config.Retention.Enabled
	app.Retention.Days = config.Retention.Days

	return &app
}
