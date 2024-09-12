package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// define a new type name envelope
type envelope map[string]interface{}

func (app *application) readIdParam(r *http.Request) (int64, error) {
	//using the ParamsFromContext to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	//Get the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id paramter")
	}
	return id, nil
}
func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//convert map to json
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	//add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
