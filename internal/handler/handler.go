package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
	"github.com/simonbrunou/parcel-tracker/internal/tracker"
)

// HandlerNotifier is implemented by the notifier package.
type HandlerNotifier interface {
	VAPIDPublicKey(ctx context.Context) (string, error)
	NotifyNewEvents(ctx context.Context, parcel model.Parcel, newEvents int)
}

type Handler struct {
	Store    store.Store
	Auth     *auth.Auth
	Tracker  *tracker.Registry
	Logger   *slog.Logger
	Notifier HandlerNotifier
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1 MB limit
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
