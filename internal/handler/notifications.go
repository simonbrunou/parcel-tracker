package handler

import (
	"net/http"

	"github.com/simonbrunou/parcel-tracker/internal/store"
)

func (h *Handler) GetVAPIDKey(w http.ResponseWriter, r *http.Request) {
	if h.Notifier == nil {
		writeError(w, http.StatusServiceUnavailable, "notifications not configured")
		return
	}
	key, err := h.Notifier.VAPIDPublicKey(r.Context())
	if err != nil || key == "" {
		writeError(w, http.StatusServiceUnavailable, "VAPID key not available")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"key": key})
}

type subscribeRequest struct {
	Endpoint string `json:"endpoint"`
	P256dh   string `json:"p256dh"`
	Auth     string `json:"auth"`
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req subscribeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Endpoint == "" || req.P256dh == "" || req.Auth == "" {
		writeError(w, http.StatusBadRequest, "endpoint, p256dh, and auth are required")
		return
	}

	sub := store.PushSubscription{
		Endpoint: req.Endpoint,
		P256dh:   req.P256dh,
		Auth:     req.Auth,
	}
	if _, err := h.Store.CreatePushSubscription(r.Context(), sub); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save subscription")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Endpoint string `json:"endpoint"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Endpoint == "" {
		writeError(w, http.StatusBadRequest, "endpoint is required")
		return
	}
	h.Store.DeletePushSubscription(r.Context(), req.Endpoint)
	w.WriteHeader(http.StatusNoContent)
}
