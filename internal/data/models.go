package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	News interface {
		Insert(news *News) error
		Get(id int64) (*News, error)
		Update(news *News) error
		Delete(id int64) error
		GetAll(title string, tags []string, author string, filters Filters) ([]*News, Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	}
	Tokens interface {
		Insert(token *Token) error
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		DeleteAllForUser(scope string, userID int64) error
	}
	Permissions interface {
		GetAllForUser(userID int64) (Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}
	Comments interface {
		Insert(comment *Comment) error
		ShowAll(id int64, filters Filters) ([]*Comment, Metadata, error)
	}
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		News:        NewsModel{pool: pool},
		Users:       UserModel{pool: pool},
		Tokens:      TokenModel{pool: pool},
		Permissions: PermissionModel{pool: pool},
		Comments:    CommentModel{pool: pool},
	}
}
