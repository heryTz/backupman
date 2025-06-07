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

var unitTest = true

func NewAppMock() *App {
	app := App{}

	var drives []drive.Drive
	var dumpers []dumper.Dumper
	var db dao.Dao
	var mailNotifier mail.MailNotifier

	if unitTest {
		drives = append(drives, &drive.DriveMock{})
		dumpers = append(dumpers, &dumper.DumperMock{})
		memoryDb := memory.NewMemoryDb()
		db = dao.Dao{
			Backup:    memory.NewBackupDaoMemory(memoryDb),
			DriveFile: memory.NewDriveFileDaoMemory(memoryDb),
			Health:    lib.MockUpHelthChecker{},
		}
		mailNotifier = mail.NewMockMailNotifier()
	} else {
		drives = append(
			drives,
			drive.NewLocalDrive("local1", "./tmp"),
			drive.NewGoogleDrive("google1", "demo", "../../service-account.json"),
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
		dbConn, err := lib.NewConnection("localhost", 3307, "root", "root", "backupman", "false")
		if err != nil {
			log.Fatal(err)
		}
		db = dao.Dao{
			Backup:    mysql.NewBackupDaoMysql(dbConn),
			DriveFile: mysql.NewDriveFileDaoMysql(dbConn),
			Health:    lib.NewHealthMysql(dbConn),
		}
		mailNotifier = mail.NewStdMailNotifier(
			"localhost",
			1026,
			"",
			"",
			"",
		)
	}

	app.Mode = APP_MODE_CLI
	app.Drives = drives
	app.Dumpers = dumpers
	app.Db = db
	app.Notifiers.Mail = mailNotifier

	return &app
}
