package main

import (
	"AwesomeProject/internal/store"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type UserWithToken struct {
	*store.User
	token string `json:"token"`
}

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		// Password: payload.Password,
	}

	// hash the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// store the user

	ctx := r.Context()
	plainToken := uuid.New().String()
	// store
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	// store the user
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerErrorHandler(w, r, err)
		}
		return
	}

	// mail

	userWithToken := UserWithToken{
		User:  user,
		token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerErrorHandler(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}
