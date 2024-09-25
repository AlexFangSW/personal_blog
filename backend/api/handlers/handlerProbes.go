package handlers

import (
	"blog/entities"
	"log/slog"
	"net/http"
)

type Probes struct{}

func NewProbes() *Probes {
	return &Probes{}
}

// ReadinessProbe
//
//	@Summary		Readiness probe
//	@Description	Readiness probe for health check
//	@Tags			healthCheck
//	@Produce		json
//	@Success		200				{object}	entities.RetSuccess[string]
//	@Router			/ready [get]
func (p *Probes) ReadinessProbe(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("ReadinessProbe")
	return entities.NewRetSuccess("ok").WriteJSON(w)
}

// LivenessProbe
//
//	@Summary		Liveness probe
//	@Description	Liveness probe for health check
//	@Tags			healthCheck
//	@Produce		json
//	@Success		200				{object}	entities.RetSuccess[string]
//	@Router			/alive [get]
func (p *Probes) LivenessProbe(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("LivenessProbe")
	return entities.NewRetSuccess("ok").WriteJSON(w)
}
