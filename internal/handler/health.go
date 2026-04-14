package handler

import "net/http"

func HealthCheck(w http.ResponseWriter, r *http.Request, h *Handler) {
	status := "ok"
	httpStatus := http.StatusOK

	if err := h.Store.Ping(r.Context()); err != nil {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	writeJSON(w, httpStatus, map[string]any{
		"status":   status,
		"carriers": h.Tracker.Available(),
	})
}
