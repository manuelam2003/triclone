package main

import (
	"errors"
	"net/http"

	"github.com/manuelam2003/triclone/internal/data"
	"github.com/manuelam2003/triclone/internal/validator"
)

func (app *application) listSettlementsHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "settled_at")
	input.Filters.SortSafelist = []string{"id", "amount", "settled_at", "payer_id", "payee_id", "-id", "-amount", "-settled_at", "-payer_id", "-payee_id"}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	settlements, metadata, err := app.models.Settlements.GetAllForGroup(groupID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"settlements": settlements, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showSettlementHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	settlementID, err := app.readIDParam(r, "settlement_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	settlement, err := app.models.Settlements.Get(settlementID, groupID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"settlement": settlement}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addSettlementHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	var input struct {
		PayerID int64   `json:"payer_id"`
		PayeeID int64   `json:"payee_id"`
		Amount  float64 `json:"amount"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	isPayerMember, err := app.models.GroupMembers.UserBelongsToGroup(input.PayerID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !isPayerMember {
		v.AddError("payer_id", "payer must be a member of the group")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	isPayeeMember, err := app.models.GroupMembers.UserBelongsToGroup(input.PayeeID, groupID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !isPayeeMember {
		v.AddError("payee_id", "payee must be a member of the group")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	settlement := &data.Settlement{
		GroupID: groupID,
		PayerID: &input.PayerID,
		PayeeID: &input.PayeeID,
		Amount:  input.Amount,
	}

	if data.ValidateSettlement(v, settlement); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Settlements.Insert(settlement)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"settlement": settlement}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSettlementHandler(w http.ResponseWriter, r *http.Request) {
	groupID, err := app.readIDParam(r, "group_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currentUser := app.contextGetUser(r)

	isMember, err := app.checkUserMembership(w, r, currentUser.ID, groupID)
	if err != nil || !isMember {
		return
	}

	settlementID, err := app.readIDParam(r, "settlement_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Settlements.Delete(groupID, settlementID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "settlement successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
