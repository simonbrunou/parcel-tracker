package handler

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

type createEventRequest struct {
	Status    model.ParcelStatus `json:"status"`
	Message   string             `json:"message"`
	Location  string             `json:"location"`
	Timestamp string             `json:"timestamp"`
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	parcelID := chi.URLParam(r, "id")
	events, err := h.Store.ListEvents(r.Context(), parcelID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	parcelID := chi.URLParam(r, "id")

	// Verify parcel exists
	if _, err := h.Store.GetParcel(r.Context(), parcelID); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "parcel not found")
		return
	}

	var req createEventRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}
	if req.Status == "" {
		req.Status = model.StatusUnknown
	}

	event := model.TrackingEvent{
		ParcelID: parcelID,
		Status:   req.Status,
		Message:  req.Message,
		Location: req.Location,
	}

	created, err := h.Store.CreateEvent(r.Context(), event)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")
	if err := h.Store.DeleteEvent(r.Context(), eventID); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "event not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete event")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
