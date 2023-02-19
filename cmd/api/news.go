package main

import (
	"fmt"
	"net/http"
	"newsportal.universityproject.net/internal/data"
	"newsportal.universityproject.net/internal/validator"
	"time"
)

func (app *application) createNewsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string    `json:"title"`
		Abstract string    `json:"abstract"`
		Tags     []string  `json:"tags"`
		Author   string    `json:"author"`
		Source   string    `json:"source"`
		Time     time.Time `json:"time"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.News{
		Title:    input.Title,
		Abstract: input.Abstract,
		Tags:     input.Tags,
		Author:   input.Author,
		Source:   input.Source,
		Time:     input.Time,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showNewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	publishTime, _ := time.Parse("2006-01-02T15:04:05Z07:00", "2023-02-18T22:09:25.000Z")
	news := data.News{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "They protested in China - and now they've gone missing",
		Abstract:  "Police have quietly gone after those who took part in the November protests that shook China",
		Tags:      []string{"China"},
		Author:    "Tessa Wong & Grace Tsoi",
		Source:    "https://www.bbc.com/news/world-asia-china-64592333",
		Time:      publishTime,
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"news": news}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
