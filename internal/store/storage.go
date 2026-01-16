package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Update(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		GetUserFeed(ctx context.Context, userId int64, feedQuery PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}
	Comments interface {
		CreateComments(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID int64, userID int64) error
		Unfollow(ctx context.Context, followerID int64, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		&PostStore{db},
		&UserStore{db},
		&CommentStore{db},
		&FollowerStore{db},
	}
}
