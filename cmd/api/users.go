package main

import (
	"AwesomeProject/internal/store"
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type userKey string

const userCtxKey userKey = "user"

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)
	// revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err := app.store.Followers.Follow(ctx, followedUser.ID, payload.UserID)
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
	// revert back to auth userID from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err := app.store.Followers.Unfollow(ctx, unfollowedUser.ID, payload.UserID)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				app.notFoundResponse(w, r, err)
			default:
				app.badRequestResponse(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtxKey).(*store.User)
	return user
}
