package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
	"github.com/simonbrunou/parcel-tracker/internal/tracker"
)

type createParcelRequest struct {
	TrackingNumber string          `json:"tracking_number"`
	Carrier        model.CarrierCode `json:"carrier"`
	Name           string          `json:"name"`
	Notes          string          `json:"notes"`
}

type updateParcelRequest struct {
	TrackingNumber string           `json:"tracking_number"`
	Carrier        model.CarrierCode `json:"carrier"`
	Name           string           `json:"name"`
	Notes          string           `json:"notes"`
	Status         model.ParcelStatus `json:"status"`
	Archived       bool             `json:"archived"`
}

func (h *Handler) ListParcels(w http.ResponseWriter, r *http.Request) {
	filter := store.ParcelFilter{
		Status: model.ParcelStatus(r.URL.Query().Get("status")),
		Search: r.URL.Query().Get("search"),
	}

	if archStr := r.URL.Query().Get("archived"); archStr != "" {
		archived := archStr == "true"
		filter.Archived = &archived
	}

	parcels, err := h.Store.ListParcels(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list parcels")
		return
	}

	writeJSON(w, http.StatusOK, parcels)
}

func (h *Handler) GetParcel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parcel, err := h.Store.GetParcel(r.Context(), id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "parcel not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get parcel")
		return
	}

	writeJSON(w, http.StatusOK, parcel)
}

func (h *Handler) CreateParcel(w http.ResponseWriter, r *http.Request) {
	var req createParcelRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TrackingNumber == "" {
		writeError(w, http.StatusBadRequest, "tracking_number is required")
		return
	}
	if req.Carrier == "" {
		req.Carrier = model.CarrierManual
	}

	parcel := model.Parcel{
		TrackingNumber: req.TrackingNumber,
		Carrier:        req.Carrier,
		Name:           req.Name,
		Notes:          req.Notes,
	}

	created, err := h.Store.CreateParcel(r.Context(), parcel)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create parcel")
		return
	}

	// Auto-refresh: fetch initial tracking events for non-manual carriers.
	if t, ok := h.Tracker.Get(created.Carrier); ok && created.Carrier != model.CarrierManual {
		if events, err := t.Track(r.Context(), created.TrackingNumber); err == nil {
			for _, e := range events {
				e.ParcelID = created.ID
				h.Store.CreateEvent(r.Context(), e)
			}
			now := time.Now().UTC()
			created.LastCheck = &now
			h.Store.UpdateParcel(r.Context(), created)
			created, _ = h.Store.GetParcel(r.Context(), created.ID)
		}
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) UpdateParcel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	existing, err := h.Store.GetParcel(r.Context(), id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "parcel not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get parcel")
		return
	}

	var req updateParcelRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing.TrackingNumber = req.TrackingNumber
	existing.Carrier = req.Carrier
	existing.Name = req.Name
	existing.Notes = req.Notes
	if req.Status != "" {
		existing.Status = req.Status
	}
	existing.Archived = req.Archived

	updated, err := h.Store.UpdateParcel(r.Context(), existing)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update parcel")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteParcel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.Store.DeleteParcel(r.Context(), id); err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "parcel not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete parcel")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshParcel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parcel, err := h.Store.GetParcel(r.Context(), id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "parcel not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get parcel")
		return
	}

	t, ok := h.Tracker.Get(parcel.Carrier)
	if !ok {
		writeError(w, http.StatusBadRequest, "unsupported carrier")
		return
	}

	events, err := t.Track(r.Context(), parcel.TrackingNumber)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "tracking failed: "+err.Error())
		return
	}

	// Load existing events to avoid duplicates.
	existing, _ := h.Store.ListEvents(r.Context(), parcel.ID)
	seen := make(map[string]bool, len(existing))
	for _, e := range existing {
		seen[tracker.EventKey(e)] = true
	}

	for _, e := range events {
		if seen[tracker.EventKey(e)] {
			continue
		}
		e.ParcelID = parcel.ID
		h.Store.CreateEvent(r.Context(), e)
	}

	// Update last_check timestamp.
	now := time.Now().UTC()
	parcel.LastCheck = &now
	h.Store.UpdateParcel(r.Context(), parcel)

	// Return updated parcel
	parcel, _ = h.Store.GetParcel(r.Context(), id)
	writeJSON(w, http.StatusOK, parcel)
}
