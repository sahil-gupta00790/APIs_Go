package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println((err))
}

//We need to send JSON -FOrmatted error msg

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//creating the json response
	env := envelope{"error": message}
	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// server error response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//logging the error
	app.logError(r, err)
	message := "The server encounted the error and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Not found error
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource is not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Not allowed method
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for tthis resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Validation errors
func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
