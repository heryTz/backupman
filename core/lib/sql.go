package lib

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewMysqlConnection(host string, port int, user string, password string, database string, tls string) (*sql.DB, error) {
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
	return dbConn, nil
}

func NewPostgresConnection(host string, port int, user string, password string, database string, tls bool) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, database)
	if tls {
		connString += "?sslmode=require"
	} else {
		connString += "?sslmode=disable"
	}
	dbConn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	return dbConn, nil
}

func NewConnection(host string, port int, user string, password string, database string, tls string) (*sql.DB, error) {
	return NewMysqlConnection(host, port, user, password, database, tls)
}

type SqlNullableTime struct {
	Time  time.Time
	Valid bool
}

func (nt *SqlNullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, true
		return nil
	}

	switch v := value.(type) {
	case []byte:
		value := string(v)
		if value == "0000-00-00 00:00:00" {
			nt.Time, nt.Valid = time.Time{}, true
			return nil
		}
		t, err := time.Parse("2006-01-02 15:04:05", value)
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
		value := string(v)
		if value == "0000-00-00 00:00:00" {
			nt.Time, nt.Valid = time.Time{}, false
			return fmt.Errorf("[SqlNonNullableTime] value is nil")
		}
		t, err := time.Parse("2006-01-02 15:04:05", value)
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
