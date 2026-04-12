package tracker

import (
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

func TestNewRegistryRegistersAllCarriers(t *testing.T) {
	t.Setenv("LAPOSTE_API_KEY", "test-key")
	r := NewRegistry()

	expectedCarriers := []model.CarrierCode{
		model.CarrierManual,
		model.CarrierMock,
		model.CarrierChronopost,
		model.CarrierLaPoste,
		model.CarrierColissimo,
		model.CarrierGLS,
		model.CarrierMondialRelay,
		model.CarrierDPD,
		model.CarrierColisPrive,
		model.CarrierRelaisColis,
		model.CarrierVintedGo,
	}

	for _, code := range expectedCarriers {
		if _, ok := r.Get(code); !ok {
			t.Errorf("expected carrier %q to be registered", code)
		}
	}
}

func TestNewRegistryHidesCarriersWithoutAPIKey(t *testing.T) {
	t.Setenv("LAPOSTE_API_KEY", "")
	r := NewRegistry()

	for _, code := range []model.CarrierCode{model.CarrierLaPoste, model.CarrierColissimo} {
		if _, ok := r.Get(code); ok {
			t.Errorf("carrier %q should not be registered without LAPOSTE_API_KEY", code)
		}
	}
}

func TestRegistryGetUnknownCarrier(t *testing.T) {
	r := NewRegistry()

	_, ok := r.Get("nonexistent")
	if ok {
		t.Error("expected ok=false for unknown carrier")
	}
}

func TestRegistryAvailable(t *testing.T) {
	t.Setenv("LAPOSTE_API_KEY", "test-key")
	r := NewRegistry()

	carriers := r.Available()
	if len(carriers) < 11 {
		t.Errorf("expected at least 11 carriers, got %d", len(carriers))
	}

	// Verify each carrier has code and name
	for _, c := range carriers {
		if c.Code == "" {
			t.Error("carrier has empty code")
		}
		if c.Name == "" {
			t.Errorf("carrier %q has empty name", c.Code)
		}
	}
}

func TestRegistryRegisterAndGet(t *testing.T) {
	r := &Registry{trackers: make(map[model.CarrierCode]Tracker)}

	r.Register(&ManualTracker{})

	got, ok := r.Get(model.CarrierManual)
	if !ok {
		t.Fatal("expected to find manual tracker")
	}
	if got.Code() != model.CarrierManual {
		t.Errorf("expected code %q, got %q", model.CarrierManual, got.Code())
	}
	if got.Name() != "Manual" {
		t.Errorf("expected name %q, got %q", "Manual", got.Name())
	}
}
