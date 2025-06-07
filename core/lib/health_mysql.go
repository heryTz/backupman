package lib

import (
	"database/sql"
	"fmt"
)

type HealthMysql struct {
	db *sql.DB
}

func NewHealthMysql(db *sql.DB) *HealthMysql {
	return &HealthMysql{db: db}
}

func (h *HealthMysql) Check() error {
	_, err := h.db.Query("SELECT 1")
	if err != nil {
		return fmt.Errorf("database health check failed => %v", err)
	}
	return nil
}
