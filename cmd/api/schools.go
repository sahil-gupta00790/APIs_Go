package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createSchoolHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "creating a new school")
}

func (app *application) showSchoolHandler(w http.ResponseWriter, r *http.Request) {
	//using the ParamsFromContext to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	//Get the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//Displayt the school id
	fmt.Fprintf(w, "show the details for school %d\n", id)

}
