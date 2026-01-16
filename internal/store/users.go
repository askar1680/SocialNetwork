package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (store *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	// up/down

	err := store.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	var (
		user      User
		createdAt time.Time
	)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := store.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &createdAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return &user, nil
}
