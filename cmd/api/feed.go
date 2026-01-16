package main

import (
	"AwesomeProject/internal/store"
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, search, filters, sort
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	posts, err := app.store.Posts.GetUserFeed(ctx, int64(101), fq)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}
