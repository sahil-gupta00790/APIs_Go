package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	//create a map to hold our healthcheck data
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	//convert map to json
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encontered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	//add a new line to make viewing on the terminal easier
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
