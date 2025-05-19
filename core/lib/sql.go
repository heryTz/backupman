package lib

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

func NewConnection(host string, port int, user string, password string, database string, tls string) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", host, port)
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.DBName = database
	if tls != "" {
		cfg.TLSConfig = tls
	}
	dbConn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	err = dbConn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}
	return dbConn, nil
}

type SqlNullableTime struct {
	Time  time.Time
	Valid bool
}

func (nt *SqlNullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		nt.Time = t
		nt.Valid = true
		return nil
	case time.Time:
		nt.Time = v
		nt.Valid = true
		return nil
	default:
		return fmt.Errorf("[SqlNullableTime] type non supporté: %T", value)
	}
}

type SqlNonNullableTime struct {
	Time  time.Time
	Valid bool
}

func (nt *SqlNonNullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return fmt.Errorf("[SqlNonNullableTime] value is nil")
	}

	switch v := value.(type) {
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		nt.Time = t
		nt.Valid = true
		return nil
	case time.Time:
		nt.Time = v
		nt.Valid = true
		return nil
	default:
		return fmt.Errorf("[SqlNonNullableTime] type non supporté: %T", value)
	}
}
