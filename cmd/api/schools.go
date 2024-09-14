package main

import (
	"errors"
	"fmt"
	"net/http"

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
	//create a school
	err = app.models.Schools.Insert(school)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//Create a location header for the newely created resource
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/schools/%d", school.ID))
	//Write the JSON Responsewith 201
	err = app.writeJson(w, http.StatusCreated, envelope{"school": school}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showSchoolHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//Fetch the specific school
	school, er := app.models.Schools.Get(id)
	if er != nil {
		switch {
		case errors.Is(er, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, er)

		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"school": school}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateSchoolhandler(w http.ResponseWriter, r *http.Request) {
	//This method does a partial replacement
	//Get the id for the school that needs updating
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//Fetch the orginal record from db
	school, er := app.models.Schools.Get(id)
	if er != nil {
		switch {
		case errors.Is(er, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, er)

		}
		return
	}
	//Create an input struct to hold the data in from the client
	//we update the input strcut to use poitners bcoz pointers have default value of nil
	//if field is nilclient did not update it
	var input struct {
		Name    *string  `json:"name"`
		Level   *string  `json:"level"`
		Contact *string  `json:"contact"`
		Phone   *string  `json:"phone"`
		Email   *string  `json:"email"`
		Website *string  `json:"website"`
		Address *string  `json:"address"`
		Mode    []string `json:"mode"`
	}
	//initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	////Copy /Update the fields/values in the school variable using the fielda in the input struct
	//school.Name = input.Name
	//school.Level = input.Level
	//school.Contact = input.Contact
	//school.Phone = input.Phone
	//school.Email = input.Email
	//school.Website = input.Website
	//school.Address = input.Address
	//school.Mode = input.Mode

	//check for updates
	if input.Name != nil {
		school.Name = *input.Name
	}
	if input.Level != nil {
		school.Level = *input.Level
	}
	if input.Contact != nil {
		school.Contact = *input.Contact
	}
	if input.Phone != nil {
		school.Phone = *input.Phone
	}
	if input.Email != nil {
		school.Email = *input.Email
	}
	if input.Website != nil {
		school.Website = *input.Website
	}
	if input.Address != nil {
		school.Address = *input.Address
	}
	if input.Address != nil {
		school.Mode = input.Mode
	}
	//initialize a new Validator instance
	v := validator.New()

	//check the map to detertmine if any validation error
	if data.ValidateSchool(v, school); !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}
	//Pass the updated school record to update emthod
	err = app.models.Schools.Update(school)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJson(w, http.StatusOK, envelope{"school": school}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteSchoolHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	//delete the school from the db
	//ssend 404 if not found
	err = app.models.Schools.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	//return a 200 status ok with a success message
	err = app.writeJson(w, http.StatusOK, envelope{"message": "school successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
