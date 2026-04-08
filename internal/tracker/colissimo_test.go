package tracker

import (
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

func TestColissimoCodeAndName(t *testing.T) {
	tracker := &ColissimoTracker{}
	if tracker.Code() != model.CarrierColissimo {
		t.Errorf("expected code %q, got %q", model.CarrierColissimo, tracker.Code())
	}
	if tracker.Name() != "Colissimo" {
		t.Errorf("expected name %q, got %q", "Colissimo", tracker.Name())
	}
}
