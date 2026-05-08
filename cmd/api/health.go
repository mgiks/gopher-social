package main

import (
	"net/http"
)

type HealthData struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	Version string `json:"version"`
}

// healthCheckHandler godoc
//
//	@Summary		Returns a healthcheck
//	@Description	Returns a healthcheck
//	@Tags			ops
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthData
//	@Failure		500	{object}	apiError
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthData{
		Status:  "ok",
		Env:     app.config.env,
		Version: apiVersion,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
