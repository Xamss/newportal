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

	router.HandlerFunc(http.MethodGet, "/v1/news", app.requirePermission("news:read", app.listNewsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/news", app.requirePermission("news:write", app.createNewsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/news/:id", app.requirePermission("news:read", app.showNewsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/news/:id", app.requirePermission("news:write", app.updateNewsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/news/:id", app.requirePermission("news:write", app.deleteNewsHandler))

	router.HandlerFunc(http.MethodGet, "/v1/news/:id/comments", app.requirePermission("news:read", app.viewCommentHandler))
	router.HandlerFunc(http.MethodPost, "/v1/news/:id/comments", app.requirePermission("news:write", app.addCommentHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
