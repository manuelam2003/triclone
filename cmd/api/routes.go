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

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.listGroupsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.createGroupHandler)
	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id", app.showGroupHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/groups/:group_id", app.updateGroupHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:group_id", app.deleteGroupHandler)

	return app.recoverPanic(router)
}
