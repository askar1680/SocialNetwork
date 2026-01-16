package main

import (
	"AwesomeProject/internal/store"
	"net/http"
)

type CreateCommentPayload struct {
	Content string `json:"content,omitempty,required"`
	UserID  int64  `json:"user_id,omitempty,required"`
	PostID  int64  `json:"post_id,omitempty,required"`
}

func (app *application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	comment := store.Comment{
		Content: payload.Content,
		PostID:  payload.PostID,
		UserID:  payload.UserID,
	}
	err := app.store.Comments.CreateComments(r.Context(), &comment)
	if err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerErrorHandler(w, r, err)
	}
}
