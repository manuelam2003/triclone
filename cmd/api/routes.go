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

	router.HandlerFunc(http.MethodGet, "/v1/users", app.listUsersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:user_id", app.showUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:user_id", app.requireActivatedUser(app.updateUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:user_id", app.requireActivatedUser(app.deleteUserHandler))

	router.HandlerFunc(http.MethodGet, "/v1/groups", app.listGroupsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id", app.showGroupHandler)
	router.HandlerFunc(http.MethodPost, "/v1/groups", app.requireActivatedUser(app.createGroupHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/groups/:group_id", app.requireActivatedUser(app.updateGroupHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:group_id", app.requireActivatedUser(app.deleteGroupHandler))

	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id/members", app.listGroupMembersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/groups/:group_id/members", app.requireActivatedUser(app.addGroupMemberHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:group_id/members/:user_id", app.removeGroupMemberHandler)
	router.HandlerFunc(http.MethodPut, "/v1/groups/:group_id/members/:user_id", app.requireActivatedUser(app.reinstateGroupMemberHandler))

	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id/expenses", app.requireActivatedUser(app.listGroupExpensesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/groups/:group_id/expenses/:expense_id", app.requireActivatedUser(app.showGroupExpenseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/groups/:group_id/expenses", app.requireActivatedUser(app.createGroupExpenseHandler))
	router.HandlerFunc(http.MethodPut, "/v1/groups/:group_id/expenses/:expense_id", app.requireActivatedUser(app.updateGroupExpenseHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:group_id/expenses/:expense_id", app.requireActivatedUser(app.deleteGroupExpenseHandler))

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
