package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "Not Found Error")
}
