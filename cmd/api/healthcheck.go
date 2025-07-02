package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	resp := `{"status":"available", "environment": %q, "version": %q}`
	resp = fmt.Sprintf(resp, app.config.env, version)

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(resp))
}
