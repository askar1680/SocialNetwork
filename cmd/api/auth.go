package main

import (
	"AwesomeProject/internal/store"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
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
		Token: plainToken,
	}

	// send email
	// TODO: email sending is not working without activating mailer accounts
	//isProdEnv := app.config.env == "production"
	//activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	//vars := struct {
	//	Username      string
	//	ActivationURL string
	//}{
	//	Username:      userWithToken.Username,
	//	ActivationURL: activationURL,
	//}
	//err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	//if err != nil {
	//	app.logger.Error("failed to send welcome email", zap.Error(err))
	//	// rollback if user creation fails
	//	if err = app.store.Users.Delete(ctx, user.ID); err != nil {
	//		app.logger.Error("failed to delete user", zap.Error(err))
	//	}
	//	app.internalServerErrorHandler(w, r, err)
	//	return
	//}

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
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerErrorHandler(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload credentials
	// fetch the user
	// generate the token -> add claims
	// send it to the client

	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerErrorHandler(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.auth.GenerateToken(claims)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}
