package main

import (
	"AwesomeProject/internal/store"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			// parse it
			// decode it
			// check it
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("empty authorization header"))
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("authorization header format must be Basic"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
				return
			}
			creds := strings.Split(string(decoded), ":")
			if len(creds) != 2 {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("invalid credential"))
				return
			}
			username := creds[0]
			password := creds[1]
			if username != app.config.auth.basic.username || password != app.config.auth.basic.password {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("invalid credential"))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, errors.New("empty authorization header"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, errors.New("authorization header format must be Bearer"))
			return
		}
		token := parts[1]
		jwtToken, err := app.auth.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
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
		app.logger.Infof("User ID: %d, Username: %s", userID, user.Username)
		ctx = context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if app.cacheStorage != nil {
		user, err := app.cacheStorage.Users.Get(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			app.logger.Infof("Getting user from Redis with ID: %d", userID)
			return user, nil
		}
	}
	dbUser, err := app.store.Users.GetByID(ctx, userID)
	if app.cacheStorage != nil {
		app.logger.Infof("Setting user to Redis with ID: %d", userID)
		_ = app.cacheStorage.Users.Set(ctx, dbUser)
	}
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		post := getPostFromContext(r)

		if post == nil {
			app.badRequestResponse(w, r, errors.New("post not found"))
			return
		}
		if user == nil {
			app.badRequestResponse(w, r, errors.New("user not found"))
			return
		}
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// check roles
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerErrorHandler(w, r, err)
			return
		}
		if !allowed {
			app.methodNotAllowedResponse(w, r, errors.New("user not allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	isAllowed := user.RoleID >= role.ID
	return isAllowed, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
