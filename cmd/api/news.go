package main

import (
	"errors"
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

	news := &data.News{
		Title:    input.Title,
		Abstract: input.Abstract,
		Tags:     input.Tags,
		Author:   input.Author,
		Source:   input.Source,
		Time:     input.Time,
	}

	v := validator.New()

	if data.ValidateNews(v, news); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.News.Insert(news)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", news.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": news}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showNewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	news, err := app.models.News.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"news": news}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateNewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	news, err := app.models.News.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title    *string    `json:"title"`
		Abstract *string    `json:"abstract"`
		Tags     []string   `json:"tags"`
		Author   *string    `json:"author"`
		Source   *string    `json:"source"`
		Time     *time.Time `json:"time"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		news.Title = *input.Title
	}
	if input.Abstract != nil {
		news.Abstract = *input.Abstract
	}
	if input.Tags != nil {
		news.Tags = input.Tags
	}
	if input.Author != nil {
		news.Author = *input.Author
	}
	if input.Source != nil {
		news.Author = *input.Author
	}
	if input.Time != nil {
		news.Time = *input.Time
	}

	v := validator.New()
	if data.ValidateNews(v, news); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.News.Update(news)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"news": news}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteNewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.News.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "news successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Tags   []string
		Author string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Author = app.readString(qs, "author", "")
	input.Tags = app.readCSV(qs, "tags", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "time", "-id", "-title", "-time"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.News.GetAll(input.Title, input.Tags, input.Author, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
