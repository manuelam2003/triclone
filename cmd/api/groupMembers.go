package main

import "net/http"

func (app *application) addGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}

func (app *application) listGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}
func (app *application) removeGroupMemberHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
}
