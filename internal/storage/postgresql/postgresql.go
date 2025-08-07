package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) BeginTransaction() (*sql.Tx, error) {
	return s.db.Begin()
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
