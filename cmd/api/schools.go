package main

import (
	"fmt"
	"net/http"
	"time"

	"Apis_go.sahil.net/internal/data" // Replace "your-package-path" with the actual package path
	"Apis_go.sahil.net/internal/validator"
)

func (app *application) createSchoolHandler(w http.ResponseWriter, r *http.Request) {
	//decoding
	var input struct {
		Name    string   `json:"name"`
		Level   string   `json:"level"`
		Contact string   `json:"contact"`
		Phone   string   `json:"phone"`
		Email   string   `json:"email"`
		Website string   `json:"website"`
		Address string   `json:"address"`
		Mode    []string `json:"mode"`
	}
	//initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//copy the values of input strcut to a new school struct
	school := &data.School{
		Name:    input.Name,
		Level:   input.Level,
		Contact: input.Website,
		Phone:   input.Phone,
		Email:   input.Email,
		Website: input.Website,
		Address: input.Address,
		Mode:    input.Mode,
	}
	//initialize a new Validator instance
	v := validator.New()

	//check the map to detertmine if any validation error
	if data.ValidateSchool(v, school); !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}
	//Display the request
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showSchoolHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//create a new instance of the school struct containing the ID we extracter from our url and some sample data
	school := data.School{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Apple Tree",
		Level:     "High School",
		Contact:   "Anna Smith",
		Phone:     "77854222",
		Address:   "14 sometghgin",
		Mode:      []string{"blended", "online"},
		Version:   1,
	}
	err = app.writeJson(w, http.StatusOK, envelope{"school": school}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
