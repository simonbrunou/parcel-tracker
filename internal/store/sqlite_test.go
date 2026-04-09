package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := NewSQLiteStore(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestCreateAndGetParcel(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, err := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC123",
		Carrier:        model.CarrierLaPoste,
		Name:           "Test Parcel",
		Notes:          "some notes",
	})
	if err != nil {
		t.Fatalf("CreateParcel: %v", err)
	}
	if p.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if p.Status != model.StatusUnknown {
		t.Errorf("expected default status %q, got %q", model.StatusUnknown, p.Status)
	}

	got, err := s.GetParcel(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetParcel: %v", err)
	}
	if got.TrackingNumber != "ABC123" {
		t.Errorf("expected tracking number %q, got %q", "ABC123", got.TrackingNumber)
	}
	if got.Name != "Test Parcel" {
		t.Errorf("expected name %q, got %q", "Test Parcel", got.Name)
	}
}

func TestGetParcelNotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	_, err := s.GetParcel(ctx, "nonexistent")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestUpdateParcel(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC123",
		Carrier:        model.CarrierLaPoste,
		Name:           "Original",
	})

	p.Name = "Updated"
	p.Status = model.StatusInTransit
	now := time.Now().UTC()
	p.LastCheck = &now

	updated, err := s.UpdateParcel(ctx, p)
	if err != nil {
		t.Fatalf("UpdateParcel: %v", err)
	}
	if updated.Name != "Updated" {
		t.Errorf("expected name %q, got %q", "Updated", updated.Name)
	}

	got, _ := s.GetParcel(ctx, p.ID)
	if got.Name != "Updated" {
		t.Errorf("expected persisted name %q, got %q", "Updated", got.Name)
	}
	if got.Status != model.StatusInTransit {
		t.Errorf("expected status %q, got %q", model.StatusInTransit, got.Status)
	}
	if got.LastCheck == nil {
		t.Error("expected LastCheck to be set")
	}
}

func TestDeleteParcel(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC123",
		Carrier:        model.CarrierManual,
	})

	if err := s.DeleteParcel(ctx, p.ID); err != nil {
		t.Fatalf("DeleteParcel: %v", err)
	}

	_, err := s.GetParcel(ctx, p.ID)
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestDeleteParcelNotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.DeleteParcel(ctx, "nonexistent")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestListParcelsFilterByStatus(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "A", Carrier: model.CarrierManual, Status: model.StatusInTransit})
	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "B", Carrier: model.CarrierManual, Status: model.StatusDelivered})
	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "C", Carrier: model.CarrierManual, Status: model.StatusInTransit})

	parcels, err := s.ListParcels(ctx, ParcelFilter{Status: model.StatusInTransit})
	if err != nil {
		t.Fatalf("ListParcels: %v", err)
	}
	if len(parcels) != 2 {
		t.Errorf("expected 2 in-transit parcels, got %d", len(parcels))
	}
}

func TestListParcelsFilterByArchived(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p1, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "A", Carrier: model.CarrierManual})
	p1.Archived = true
	s.UpdateParcel(ctx, p1)
	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "B", Carrier: model.CarrierManual})

	notArchived := false
	parcels, err := s.ListParcels(ctx, ParcelFilter{Archived: &notArchived})
	if err != nil {
		t.Fatalf("ListParcels: %v", err)
	}
	if len(parcels) != 1 {
		t.Errorf("expected 1 non-archived parcel, got %d", len(parcels))
	}
	if parcels[0].TrackingNumber != "B" {
		t.Errorf("expected tracking number B, got %s", parcels[0].TrackingNumber)
	}

	archived := true
	parcels, err = s.ListParcels(ctx, ParcelFilter{Archived: &archived})
	if err != nil {
		t.Fatalf("ListParcels: %v", err)
	}
	if len(parcels) != 1 {
		t.Errorf("expected 1 archived parcel, got %d", len(parcels))
	}
}

func TestListParcelsFilterBySearch(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "TRACK123", Carrier: model.CarrierManual, Name: "My Phone"})
	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "OTHER456", Carrier: model.CarrierManual, Name: "My Laptop"})

	// Search by tracking number
	parcels, _ := s.ListParcels(ctx, ParcelFilter{Search: "TRACK"})
	if len(parcels) != 1 {
		t.Errorf("expected 1 result searching by tracking number, got %d", len(parcels))
	}

	// Search by name
	parcels, _ = s.ListParcels(ctx, ParcelFilter{Search: "Laptop"})
	if len(parcels) != 1 {
		t.Errorf("expected 1 result searching by name, got %d", len(parcels))
	}

	// Search matching both
	parcels, _ = s.ListParcels(ctx, ParcelFilter{Search: "My"})
	if len(parcels) != 2 {
		t.Errorf("expected 2 results for broad search, got %d", len(parcels))
	}
}

func TestListParcelsEmpty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	parcels, err := s.ListParcels(ctx, ParcelFilter{})
	if err != nil {
		t.Fatalf("ListParcels: %v", err)
	}
	if parcels == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(parcels) != 0 {
		t.Errorf("expected 0 parcels, got %d", len(parcels))
	}
}

func TestListParcelsOrder(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p1, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "FIRST", Carrier: model.CarrierManual})
	time.Sleep(10 * time.Millisecond)
	s.CreateParcel(ctx, model.Parcel{TrackingNumber: "SECOND", Carrier: model.CarrierManual})

	parcels, _ := s.ListParcels(ctx, ParcelFilter{})
	if len(parcels) != 2 {
		t.Fatalf("expected 2 parcels, got %d", len(parcels))
	}
	// Most recently updated first
	if parcels[0].TrackingNumber != "SECOND" {
		t.Errorf("expected SECOND first (newest), got %s", parcels[0].TrackingNumber)
	}

	// Update the first parcel to make it newest
	p1.Name = "updated"
	s.UpdateParcel(ctx, p1)

	parcels, _ = s.ListParcels(ctx, ParcelFilter{})
	if parcels[0].TrackingNumber != "FIRST" {
		t.Errorf("expected FIRST first after update, got %s", parcels[0].TrackingNumber)
	}
}

func TestCreateAndListEvents(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	e1, err := s.CreateEvent(ctx, model.TrackingEvent{
		ParcelID: p.ID,
		Status:   model.StatusInTransit,
		Message:  "In transit",
	})
	if err != nil {
		t.Fatalf("CreateEvent: %v", err)
	}
	if e1.ID == "" {
		t.Fatal("expected non-empty event ID")
	}

	events, err := s.ListEvents(ctx, p.ID)
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Message != "In transit" {
		t.Errorf("expected message %q, got %q", "In transit", events[0].Message)
	}
}

func TestCreateEventUpdatesParcelStatus(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	if p.Status != model.StatusUnknown {
		t.Fatalf("expected initial status %q, got %q", model.StatusUnknown, p.Status)
	}

	s.CreateEvent(ctx, model.TrackingEvent{
		ParcelID: p.ID,
		Status:   model.StatusDelivered,
		Message:  "Delivered",
	})

	updated, _ := s.GetParcel(ctx, p.ID)
	if updated.Status != model.StatusDelivered {
		t.Errorf("expected parcel status updated to %q, got %q", model.StatusDelivered, updated.Status)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	e, _ := s.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "test"})

	if err := s.DeleteEvent(ctx, e.ID); err != nil {
		t.Fatalf("DeleteEvent: %v", err)
	}

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected 0 events after delete, got %d", len(events))
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	err := s.DeleteEvent(ctx, "nonexistent")
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestDeleteParcelCascadesEvents(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	s.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "event1"})
	s.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "event2"})

	s.DeleteParcel(ctx, p.ID)

	events, _ := s.ListEvents(ctx, p.ID)
	if len(events) != 0 {
		t.Errorf("expected events to be cascade-deleted, got %d", len(events))
	}
}

func TestSettings(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	// Get missing key returns empty string
	val, err := s.GetSetting(ctx, "missing_key")
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string for missing key, got %q", val)
	}

	// Set and get
	if err := s.SetSetting(ctx, "test_key", "test_value"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	val, err = s.GetSetting(ctx, "test_key")
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if val != "test_value" {
		t.Errorf("expected %q, got %q", "test_value", val)
	}

	// Upsert
	if err := s.SetSetting(ctx, "test_key", "updated"); err != nil {
		t.Fatalf("SetSetting upsert: %v", err)
	}
	val, _ = s.GetSetting(ctx, "test_key")
	if val != "updated" {
		t.Errorf("expected %q after upsert, got %q", "updated", val)
	}
}

func TestCreateEventDefaultTimestamp(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	e, err := s.CreateEvent(ctx, model.TrackingEvent{
		ParcelID: p.ID,
		Message:  "test",
	})
	if err != nil {
		t.Fatalf("CreateEvent: %v", err)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero default timestamp")
	}
}

func TestCreateEventWithExplicitTimestamp(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	e, err := s.CreateEvent(ctx, model.TrackingEvent{
		ParcelID:  p.ID,
		Message:   "test",
		Timestamp: ts,
	})
	if err != nil {
		t.Fatalf("CreateEvent: %v", err)
	}
	if !e.Timestamp.Equal(ts) {
		t.Errorf("expected timestamp %v, got %v", ts, e.Timestamp)
	}
}

func TestListEventsEmpty(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	events, err := s.ListEvents(ctx, "nonexistent-parcel")
	if err != nil {
		t.Fatalf("ListEvents: %v", err)
	}
	if events == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestParcelArchivedBoolConversion(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	p, _ := s.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	if p.Archived {
		t.Error("expected Archived=false by default")
	}

	p.Archived = true
	s.UpdateParcel(ctx, p)

	got, _ := s.GetParcel(ctx, p.ID)
	if !got.Archived {
		t.Error("expected Archived=true after update")
	}
}
