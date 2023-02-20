package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createNewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showNewsHandler)
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateNewsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteNewsHandler)
	return router
}
