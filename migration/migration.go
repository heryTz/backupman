package migration

import (
	"fmt"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/migration/mysql"
)

func Run(db application.DbConfig) error {
	switch db.(type) {
	case application.MysqlDbConfig:
		return mysql.RunMysql(db.(application.MysqlDbConfig))
	case application.MemoryDbConfig:
		return nil
	default:
		return fmt.Errorf("unsupported database type for migration")
	}
}
