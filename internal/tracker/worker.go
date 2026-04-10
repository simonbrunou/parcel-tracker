package tracker

import (
	"context"
	"log/slog"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
)

// Worker periodically refreshes tracking data for active parcels.
type Worker struct {
	Store    store.Store
	Registry *Registry
	Interval time.Duration
	Logger   *slog.Logger
}

// Run starts the background refresh loop. It blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) {
	w.Logger.Info("tracking worker started", "interval", w.Interval)

	// Run once immediately on startup.
	w.refreshAll(ctx)

	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.Logger.Info("tracking worker stopped")
			return
		case <-ticker.C:
			w.refreshAll(ctx)
		}
	}
}

func (w *Worker) refreshAll(ctx context.Context) {
	parcels, err := w.Store.ListActiveParcels(ctx)
	if err != nil {
		w.Logger.Error("worker: failed to list parcels", "error", err)
		return
	}

	for _, p := range parcels {
		if ctx.Err() != nil {
			return
		}
		w.refreshParcel(ctx, p)
	}
}

func (w *Worker) refreshParcel(ctx context.Context, p model.Parcel) {
	t, ok := w.Registry.Get(p.Carrier)
	if !ok {
		return
	}

	result, err := t.Track(ctx, p.TrackingNumber)
	if err != nil {
		w.Logger.Warn("worker: tracking failed",
			"parcel_id", p.ID,
			"carrier", p.Carrier,
			"tracking_number", p.TrackingNumber,
			"error", err,
		)
		return
	}

	// Load existing events to avoid duplicates.
	existing, err := w.Store.ListEvents(ctx, p.ID)
	if err != nil {
		w.Logger.Error("worker: failed to list events", "parcel_id", p.ID, "error", err)
		return
	}
	seen := make(map[string]bool, len(existing))
	for _, e := range existing {
		seen[EventKey(e)] = true
	}

	newCount := 0
	for _, e := range result.Events {
		if seen[EventKey(e)] {
			continue
		}
		e.ParcelID = p.ID
		if _, err := w.Store.CreateEvent(ctx, e); err != nil {
			w.Logger.Error("worker: failed to create event", "parcel_id", p.ID, "error", err)
			continue
		}
		newCount++
	}

	// Re-read parcel to pick up status changes made by CreateEvent,
	// then update last_check. Without this, UpdateParcel would overwrite
	// the status back to its pre-refresh value.
	updated, err := w.Store.GetParcel(ctx, p.ID)
	if err != nil {
		w.Logger.Error("worker: failed to re-read parcel", "parcel_id", p.ID, "error", err)
		return
	}
	now := time.Now().UTC()
	updated.LastCheck = &now
	updated.EstimatedDelivery = result.EstimatedDelivery
	if _, err := w.Store.UpdateParcel(ctx, updated); err != nil {
		w.Logger.Error("worker: failed to update parcel", "parcel_id", p.ID, "error", err)
	}

	if newCount > 0 {
		w.Logger.Info("worker: new events",
			"parcel_id", p.ID,
			"carrier", p.Carrier,
			"new_events", newCount,
		)
	}
}

// EventKey returns a deduplication key for a tracking event.
// Used by both the background worker and the HTTP refresh handler.
func EventKey(e model.TrackingEvent) string {
	return e.Timestamp.UTC().Format(time.RFC3339) + "|" + string(e.Status) + "|" + e.Message
}
