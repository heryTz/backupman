package core

import (
	"github.com/herytz/backupman/core/dao"
	"github.com/herytz/backupman/core/dao/memory"
	"github.com/herytz/backupman/core/drive"
	"github.com/herytz/backupman/core/dumper"
)

func NewAppMock() *App {
	app := App{}

	drives := make([]drive.Drive, 0, 1)
	drives = append(drives, &drive.DriveMock{})

	dumpers := make([]dumper.Dumper, 0, 1)
	dumpers = append(dumpers, &dumper.DumperMock{})

	memoryDb := memory.NewMemoryDb()
	db := dao.Dao{
		Backup:    memory.NewBackupDaoMemory(memoryDb),
		DriveFile: memory.NewDriveFileDaoMemory(memoryDb),
	}

	// dumpers := make([]dumper.Dumper, 0, 1)
	// dumpers = append(dumpers, dumper.NewMysqlDumper("mysql1", "./tmp", "localhost", 3306, "root", "root", "backupman", "false"))

	// dbConn, err := lib.NewConnection("localhost", 3306, "root", "root", "backupman", "false")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db := dao.Dao{
	// 	Backup:    mysql.NewBackupDaoMysql(dbConn),
	// 	DriveFile: mysql.NewDriveFileDaoMysql(dbConn),
	// }

	app.Config.AppUrl = "http://localhost:8080"
	app.Config.BackupCron = "0 0 * * *"
	app.ApiKeys = []string{"123"}
	app.Drives = drives
	app.Dumpers = dumpers
	app.Db = db

	return &app
}
