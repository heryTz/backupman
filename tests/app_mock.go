package tests

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/dao"
	"github.com/herytz/backupman/core/dao/memory"
	"github.com/herytz/backupman/core/dao/mysql"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/dumper"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/mailer"
	"github.com/herytz/backupman/core/notifier"
	"github.com/herytz/backupman/migration"
)

var unitTest = true

func NewAppMock() *application.App {
	app := application.App{}

	var drives []drive.Drive
	var dumpers []dumper.Dumper
	var db dao.Dao
	var notifiers []notifier.Notifier

	if unitTest {
		drives = append(drives, &drive.DriveMock{})
		dumpers = append(dumpers, &dumper.DumperMock{})
		memoryDb := memory.NewMemoryDb()
		db = dao.Dao{
			Backup:    memory.NewBackupDaoMemory(memoryDb),
			DriveFile: memory.NewDriveFileDaoMemory(memoryDb),
			Health:    lib.MockUpHelthChecker{},
		}
		notifiers = append(notifiers, &notifier.MockNotifier{})
	} else {
		drives = append(
			drives,
			drive.NewLocalDrive("local1", "./tmp"),
			drive.NewGoogleDrive("google1", "backupman", "../google-client-secret.json", "../google-token.json"),
		)
		dumpers = append(dumpers, dumper.NewMysqlDumper(
			"mysql1",
			"./tmp",
			"localhost",
			3307,
			"root",
			"root",
			"backupman",
			"false",
		))
		err := migration.Run(application.MysqlDbConfig{
			Host:     "localhost",
			Port:     3307,
			User:     "root",
			Password: "root",
			Database: "backupman",
			Tls:      "false",
		})
		if err != nil {
			log.Fatal(err)
		}
		dbConn, err := lib.NewConnection("localhost", 3307, "root", "root", "backupman", "false")
		if err != nil {
			log.Fatal(err)
		}
		db = dao.Dao{
			Backup:    mysql.NewBackupDaoMysql(dbConn),
			DriveFile: mysql.NewDriveFileDaoMysql(dbConn),
			Health:    lib.NewHealthMysql(dbConn),
		}
		mailerTransport := mailer.NewStdMailer(
			"localhost",
			1026,
			"",
			"",
			"",
		)
		destinations := []mailer.Recipient{}
		destinations = append(destinations, mailer.Recipient{Name: "John Doe", Email: "john.doe@yopmail.com"})
		notifiers = append(notifiers, notifier.NewMailNotifier(mailerTransport, db, destinations))
	}

	app.Mode = application.APP_MODE_CLI
	app.Drives = drives
	app.Dumpers = dumpers
	app.Db = db
	app.Notifiers = notifiers

	return &app
}
