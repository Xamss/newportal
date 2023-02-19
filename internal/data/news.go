package data

import (
	"newsportal.universityproject.net/internal/validator"
	"time"
)

type News struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
	Tags      []string  `json:"tags,omitempty"`
	Author    string    `json:"author,omitempty"`
	Source    string    `json:"source,omitempty"`
	Time      time.Time `json:"year,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, news *News) {
	v.Check(news.Title != "", "title", "must be provided")
	v.Check(len(news.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(news.Time.Compare(time.Now()) <= 0, "time", "posted time cannot be in the future")
	v.Check(news.Tags != nil, "tags", "must be provided")
	v.Check(len(news.Tags) >= 1, "tags", "must contain at least 1 tag")
	v.Check(len(news.Tags) <= 5, "tags", "must not contain more than 5 tags")
	v.Check(validator.Unique(news.Tags), "tags", "must not contain duplicate values")
	v.Check(news.Abstract != "", "abstract", "must be provided")
	v.Check(len(news.Title) <= 1500, "abstract", "must not be more than 1500 bytes long")
	v.Check(news.Author != "", "author", "must be provided")
	v.Check(news.Author != "", "author", "must be less than 500 bytes long. Please, use shorter version if possible")
	v.Check(news.Source != "", "source", "must refer to the original post")
}
