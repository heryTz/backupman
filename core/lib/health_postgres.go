package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthPostgres struct {
	db *pgxpool.Pool
}

func NewHealthPostgres(db *pgxpool.Pool) *HealthPostgres {
	return &HealthPostgres{db: db}
}

func (h *HealthPostgres) Check() error {
	_, err := h.db.Query(context.Background(), "SELECT 1")
	if err != nil {
		return fmt.Errorf("database health check failed => %v", err)
	}
	return nil
}
