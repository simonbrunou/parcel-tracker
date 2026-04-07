package handler

import "net/http"

func HealthCheck(w http.ResponseWriter, r *http.Request, h *Handler) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"carriers": h.Tracker.Available(),
	})
}
