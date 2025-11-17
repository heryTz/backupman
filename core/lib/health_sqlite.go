package lib

import (
	"database/sql"
	"fmt"
)

type HealthSqlite struct {
	db *sql.DB
}

func NewHealthSqlite(db *sql.DB) *HealthSqlite {
	return &HealthSqlite{db: db}
}

func (h *HealthSqlite) Check() error {
	err := h.db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}
	return nil
}
