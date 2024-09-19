package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.requireActivatedUser(app.listGroupsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.requireActivatedUser(app.createGroupHandler))
	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id", app.showGroupHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/groups/:group_id", app.updateGroupHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:group_id", app.deleteGroupHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
