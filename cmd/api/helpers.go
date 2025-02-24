package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"Apis_go.sahil.net/internal/validator"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	//Use http.MaxByteReader() to limit the size of the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	//decode request body into the target destination
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unMarshalTypeError *json.UnmarshalTypeError //difference between what we expect and what we get
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly found JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF): //syntax problem with json as it is being decoded
			return errors.New("body contains badly formed JSON")
		//checkinf for wrong types
		case errors.As(err, &unMarshalTypeError):
			if unMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unMarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			//Unmapable fields
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
			//too large
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		//pass non-nil pointer error

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err

		}

	}
	//check iof some value is given twice
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}
	return nil
}

// readstring()-method returns a string value from a querry parameters string or it returns a default value if no mathcing key is found
func (app *application) readstring(qs url.Values, key string, defaultValue string) string {
	//get the value
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	return value

}

// The readCSV() method splits a value into a slice based on the comma seprator
// IF no matching key is found then the default value is returned.
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	//split the string based on comma
	return strings.Split(value, ",")
}

// readINT() method convcerts a string value from query to an integer
// if value cannot be converted to interger , then validation error is added to validation map errors
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	//get the value
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return intValue

}
