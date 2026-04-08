package tracker

import (
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testMondialRelayHTML = `<!DOCTYPE html>
<html>
<body>
<div class="infos-account">
  <div><p>01/06/2025</p></div>
  <div>
    <div><div><p>10:00</p><p>Colis pris en charge par Mondial Relay</p></div></div>
    <div><div><p>14:30</p><p>Colis en transit vers le hub régional</p></div></div>
  </div>
</div>
<div class="infos-account">
  <div><p>02/06/2025</p></div>
  <div>
    <div><div><p>08:00</p><p>Colis arrivé au Point Relais</p></div></div>
    <div><div><p>16:45</p><p>Colis livré au Point Relais</p></div></div>
  </div>
</div>
</body>
</html>`

func TestParseMondialRelayHTML(t *testing.T) {
	events, err := parseMondialRelayHTML([]byte(testMondialRelayHTML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}

	tests := []struct {
		index  int
		status model.ParcelStatus
		msg    string
	}{
		{0, model.StatusInfoReceived, "Colis pris en charge par Mondial Relay"},
		{1, model.StatusInTransit, "Colis en transit vers le hub régional"},
		{2, model.StatusInTransit, "Colis arrivé au Point Relais"},
		{3, model.StatusDelivered, "Colis livré au Point Relais"},
	}

	for _, tt := range tests {
		e := events[tt.index]
		if e.Status != tt.status {
			t.Errorf("event[%d]: expected status %q, got %q", tt.index, tt.status, e.Status)
		}
		if e.Message != tt.msg {
			t.Errorf("event[%d]: expected message %q, got %q", tt.index, tt.msg, e.Message)
		}
	}
}

func TestParseMondialRelayHTMLEmpty(t *testing.T) {
	html := `<!DOCTYPE html><html><body><p>No tracking data</p></body></html>`
	events, err := parseMondialRelayHTML([]byte(html))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestMondialRelayCodeAndName(t *testing.T) {
	tracker := &MondialRelayTracker{}
	if tracker.Code() != model.CarrierMondialRelay {
		t.Errorf("expected code %q, got %q", model.CarrierMondialRelay, tracker.Code())
	}
	if tracker.Name() != "Mondial Relay" {
		t.Errorf("expected name %q, got %q", "Mondial Relay", tracker.Name())
	}
}

func TestMondialRelayTrackingNumberFormat(t *testing.T) {
	tracker := &MondialRelayTracker{}
	_, err := tracker.Track(t.Context(), "12345678")
	if err == nil {
		t.Fatal("expected error for tracking number without postal code")
	}
}

func TestMapMondialRelayStatus(t *testing.T) {
	tests := []struct {
		label  string
		status model.ParcelStatus
	}{
		{"Colis livré", model.StatusDelivered},
		{"Colis retiré", model.StatusDelivered},
		{"Colis distribué", model.StatusDelivered},
		{"Colis en cours de livraison", model.StatusOutForDelivery},
		{"Colis pris en charge", model.StatusInfoReceived},
		{"Colis déposé en Point Relais", model.StatusInfoReceived},
		{"Colis en transit", model.StatusInTransit},
		{"Anomalie de traitement", model.StatusFailed},
	}

	for _, tt := range tests {
		got := mapMondialRelayStatus(tt.label)
		if got != tt.status {
			t.Errorf("mapMondialRelayStatus(%q) = %q, want %q", tt.label, got, tt.status)
		}
	}
}
