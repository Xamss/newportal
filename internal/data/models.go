package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	News interface {
		Insert(movie *News) error
		Get(id int64) (*News, error)
		Update(movie *News) error
		Delete(id int64) error
	}
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		News: NewsModel{pool: pool},
	}
}
