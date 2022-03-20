package database

import (
	"carWash/internal/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.DBName,
		cfg.Postgres.Password,
		cfg.Postgres.SSLMode,
	),
	)

	if err != nil {
		return nil, fmt.Errorf("database.NewPostgresDB: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("database.NewPostgresDB: %w", err)
	}
	return conn, nil

}
