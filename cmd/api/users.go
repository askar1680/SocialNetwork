package main

import (
	"AwesomeProject/internal/store"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || userID < 1 {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	user, err := app.getUser(ctx, userID)
	if err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}
	if user == nil {
		app.unauthorizedErrorResponse(w, r, errors.New("user not found in middleware"))
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Followers.Follow(ctx, followedUser.ID, followedID)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromContext(r)
	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err = app.store.Followers.Unfollow(ctx, unfollowedUser.ID, unfollowedID)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}
