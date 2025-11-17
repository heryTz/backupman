package migration

import (
	"fmt"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/migration/mysql"
	"github.com/herytz/backupman/migration/postgres"
	"github.com/herytz/backupman/migration/sqlite"
)

func Run(db application.DbConfig) error {
	switch db.(type) {
	case application.MysqlDbConfig:
		return mysql.RunMysql(db.(application.MysqlDbConfig))
	case application.PostgresDbConfig:
		return postgres.RunPostgres(db.(application.PostgresDbConfig))
	case application.SqliteDbConfig:
		return sqlite.RunSqlite(db.(application.SqliteDbConfig))
	case application.MemoryDbConfig:
		return nil
	default:
		return fmt.Errorf("unsupported database type for migration")
	}
}
