package store

import (
	"context"
	"database/sql"
)

type Role struct {
	ID          int64  `json:"id"`
	Level       int64  `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := "SELECT id, level, name, description FROM roles WHERE name = $1"
	row := s.db.QueryRowContext(ctx, query, roleName)
	var r Role
	if err := row.Scan(&r.ID, &r.Level, &r.Name, &r.Description); err != nil {
		return nil, err
	}
	return &r, nil
}
