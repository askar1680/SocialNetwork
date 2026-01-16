package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

var (
	ErrorNotFound        = errors.New("Resource not found")
	QueryTimeOutDuration = time.Second * 5
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int64     `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (store *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := store.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (store *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version from posts where id = $1
    `
	var (
		post      Post
		createdAt time.Time
		updatedAt time.Time
	)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := store.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&createdAt,
		&updatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	post.CreatedAt = createdAt.Format(time.RFC3339)
	post.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &post, nil
}

func (store *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts SET title = $2, content = $3, tags = $4, version = version + 1, updated_at = NOW()
		WHERE id = $1 AND version = $5
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := store.db.QueryRowContext(ctx, query, post.ID, post.Title, post.Content, pq.Array(post.Tags), post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (store *PostStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := store.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorNotFound
	}
	return err
}

func (store *PostStore) GetUserFeed(ctx context.Context, userId int64, feedQuery PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT
			posts.id,
			posts.user_id,
			posts.title,
			posts.content,
			posts.created_at,
			posts.tags,
			users.username,
			COUNT(comments.id) AS comments_count from posts
		LEFT JOIN comments ON posts.id = comments.post_id
		LEFT JOIN users ON posts.user_id = users.id
    	JOIN followers ON posts.user_id = followers.follower_id OR posts.user_id = $1
    	WHERE 
    	    followers.user_id = $1 AND
    		(posts.title ILIKE '%' || $4 || '%' OR posts.content ILIKE '%' || $4 || '%') AND
    		(posts.tags @> $5 OR $5 = '{}')
    	GROUP BY posts.id, users.username
    	ORDER BY posts.created_at ` + feedQuery.Sort + `
    	LIMIT $2 OFFSET $3
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows, err := store.db.QueryContext(ctx, query, userId, feedQuery.Limit, feedQuery.Offset, feedQuery.Search, pq.Array(feedQuery.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var postsWithMetadata []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		postsWithMetadata = append(postsWithMetadata, p)
	}
	return postsWithMetadata, nil
}
