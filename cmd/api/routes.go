package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/news", app.createNewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/news/:id", app.showNewsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/news/:id", app.updateNewsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/news/:id", app.deleteNewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/news", app.listMoviesHandler)

	return app.recoverPanic(app.rateLimit(router))
}
