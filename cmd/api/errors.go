package main

import (
	"net/http"
)

func (app *application) internalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal Server Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad Request Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infof("Not Found Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "Not Found Error")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infof("Unauthorized Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized Error")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infof("Unauthorized Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized Error")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Infof("Method Not Allowed Error %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed Error")
}
