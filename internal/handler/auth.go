package handler

import (
	"net/http"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
)

type loginRequest struct {
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.Auth.Login(r.Context(), req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	auth.SetSessionCookie(w, token)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	auth.ClearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	token := auth.ExtractToken(r)
	if token == "" {
		writeJSON(w, http.StatusOK, map[string]bool{"authenticated": false, "configured": h.Auth.IsConfigured(r.Context())})
		return
	}

	if err := h.Auth.Verify(r.Context(), token); err != nil {
		writeJSON(w, http.StatusOK, map[string]bool{"authenticated": false, "configured": true})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"authenticated": true, "configured": true})
}

type setupRequest struct {
	Password string `json:"password"`
}

func (h *Handler) Setup(w http.ResponseWriter, r *http.Request) {
	if h.Auth.IsConfigured(r.Context()) {
		writeError(w, http.StatusConflict, "already configured")
		return
	}

	var req setupRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Password) < 4 {
		writeError(w, http.StatusBadRequest, "password must be at least 4 characters")
		return
	}

	if err := h.Auth.Setup(r.Context(), req.Password); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to setup")
		return
	}

	token, err := h.Auth.Login(r.Context(), req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to login after setup")
		return
	}

	auth.SetSessionCookie(w, token)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
