package tracker

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
)

func newTestWorkerStore(t *testing.T) store.Store {
	t.Helper()
	s, err := store.NewSQLiteStore(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestEventKey(t *testing.T) {
	ts := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
	e := model.TrackingEvent{
		Status:    model.StatusInTransit,
		Message:   "In transit",
		Timestamp: ts,
	}

	key := EventKey(e)
	expected := "2025-06-01T10:00:00Z|in_transit|In transit"
	if key != expected {
		t.Errorf("expected key %q, got %q", expected, key)
	}
}

func TestEventKeyDifferentEvents(t *testing.T) {
	ts := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)

	e1 := model.TrackingEvent{Status: model.StatusInTransit, Message: "In transit", Timestamp: ts}
	e2 := model.TrackingEvent{Status: model.StatusDelivered, Message: "Delivered", Timestamp: ts}

	if EventKey(e1) == EventKey(e2) {
		t.Error("different events should have different keys")
	}
}

func TestEventKeySameEventSameKey(t *testing.T) {
	ts := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)

	e1 := model.TrackingEvent{Status: model.StatusInTransit, Message: "In transit", Timestamp: ts}
	e2 := model.TrackingEvent{Status: model.StatusInTransit, Message: "In transit", Timestamp: ts}

	if EventKey(e1) != EventKey(e2) {
		t.Error("same events should have same keys")
	}
}

func TestRefreshAllSkipsArchivedParcels(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	// Create an archived parcel with mock carrier
	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ARCHIVED1",
		Carrier:        model.CarrierMock,
		Status:         model.StatusInTransit,
	})
	p.Archived = true
	s.UpdateParcel(ctx, p)

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}
	reg.Register(&MockTracker{})

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	w.refreshAll(ctx)

	// Archived parcel should not have been refreshed (no events)
	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected 0 events for archived parcel, got %d", len(events))
	}
}

func TestRefreshAllSkipsTerminalStatuses(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	// Create a delivered parcel
	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "DELIVERED1",
		Carrier:        model.CarrierMock,
		Status:         model.StatusDelivered,
	})

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}
	reg.Register(&MockTracker{})

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	w.refreshAll(ctx)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected 0 events for delivered parcel, got %d", len(events))
	}
}

func TestRefreshAllSkipsManualCarrier(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "MANUAL1",
		Carrier:        model.CarrierManual,
		Status:         model.StatusInTransit,
	})

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}
	reg.Register(&ManualTracker{})

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	w.refreshAll(ctx)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected 0 events for manual parcel, got %d", len(events))
	}
}

func TestRefreshParcelCreatesEvents(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "MOCK1",
		Carrier:        model.CarrierMock,
		Status:         model.StatusUnknown,
	})

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}
	reg.Register(&MockTracker{})

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	w.refreshParcel(ctx, p)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 5 {
		t.Fatalf("expected 5 events from mock tracker, got %d", len(events))
	}

	// Verify parcel status was updated
	updated, _ := s.GetParcel(ctx, p.ID)
	if updated.Status != model.StatusDelivered {
		t.Errorf("expected status %q, got %q", model.StatusDelivered, updated.Status)
	}
	if updated.LastCheck == nil {
		t.Error("expected LastCheck to be set")
	}
}

func TestRefreshParcelDeduplicatesEvents(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "MOCK1",
		Carrier:        model.CarrierMock,
		Status:         model.StatusUnknown,
	})

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}
	reg.Register(&MockTracker{})

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	// Refresh twice
	w.refreshParcel(ctx, p)
	p, _ = s.GetParcel(ctx, p.ID) // Re-read to get updated state
	w.refreshParcel(ctx, p)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 5 {
		t.Errorf("expected 5 events after two refreshes (dedup), got %d", len(events))
	}
}

func TestRefreshParcelUnknownCarrier(t *testing.T) {
	s := newTestWorkerStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "UNKNOWN1",
		Carrier:        "nonexistent",
		Status:         model.StatusUnknown,
	})

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	// Should not panic
	w.refreshParcel(ctx, p)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected 0 events for unknown carrier, got %d", len(events))
	}
}

func TestWorkerRunStopsOnContextCancel(t *testing.T) {
	s := newTestWorkerStore(t)

	reg := &Registry{trackers: make(map[model.CarrierCode]Tracker)}

	w := &Worker{
		Store:    s,
		Registry: reg,
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// Worker stopped as expected
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not stop after context cancel")
	}
}
