package data

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
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
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		News: NewsModel{pool: pool},
	}
}
