package tracker

import (
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testColisPriveHTML = `<!DOCTYPE html>
<html>
<body>
<div class="delivery-info"><span>Livraison prévue le 03/06/2025</span></div>
<table>
  <tr class="bandeauText">
    <td>01/06/2025 10:00</td>
    <td>Colis réceptionné par Colis Privé</td>
  </tr>
  <tr class="bandeauText">
    <td>02/06/2025 08:30</td>
    <td>Colis en cours d'acheminement</td>
  </tr>
  <tr class="bandeauText">
    <td>03/06/2025 07:00</td>
    <td>Colis chargé pour livraison</td>
  </tr>
  <tr class="bandeauText">
    <td>03/06/2025 14:23</td>
    <td>Colis livré</td>
  </tr>
</table>
</body>
</html>`

func TestParseColisPriveHTML(t *testing.T) {
	result, err := parseColisPriveHTML([]byte(testColisPriveHTML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(result.Events))
	}

	if result.EstimatedDelivery == nil {
		t.Fatal("expected estimated delivery to be set")
	}

	tests := []struct {
		index  int
		status model.ParcelStatus
		msg    string
	}{
		{0, model.StatusInfoReceived, "Colis réceptionné par Colis Privé"},
		{1, model.StatusInTransit, "Colis en cours d'acheminement"},
		{2, model.StatusOutForDelivery, "Colis chargé pour livraison"},
		{3, model.StatusDelivered, "Colis livré"},
	}

	for _, tt := range tests {
		e := result.Events[tt.index]
		if e.Status != tt.status {
			t.Errorf("event[%d]: expected status %q, got %q", tt.index, tt.status, e.Status)
		}
		if e.Message != tt.msg {
			t.Errorf("event[%d]: expected message %q, got %q", tt.index, tt.msg, e.Message)
		}
	}
}

func TestParseColisPriveHTMLEmpty(t *testing.T) {
	html := `<!DOCTYPE html><html><body><p>No data</p></body></html>`
	result, err := parseColisPriveHTML([]byte(html))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(result.Events))
	}
	if result.EstimatedDelivery != nil {
		t.Error("expected nil estimated delivery for empty response")
	}
}

func TestColisPriveCodeAndName(t *testing.T) {
	tracker := &ColisPriveTracker{}
	if tracker.Code() != model.CarrierColisPrive {
		t.Errorf("expected code %q, got %q", model.CarrierColisPrive, tracker.Code())
	}
	if tracker.Name() != "Colis Privé" {
		t.Errorf("expected name %q, got %q", "Colis Privé", tracker.Name())
	}
}

func TestMapColisPriveStatus(t *testing.T) {
	tests := []struct {
		label  string
		status model.ParcelStatus
	}{
		{"Colis livré", model.StatusDelivered},
		{"Colis remis au destinataire", model.StatusDelivered},
		{"Colis en cours de livraison", model.StatusOutForDelivery},
		{"Colis chargé pour livraison", model.StatusOutForDelivery},
		{"Colis réceptionné", model.StatusInfoReceived},
		{"Colis pris en charge", model.StatusInfoReceived},
		{"Colis en cours d'acheminement", model.StatusInTransit},
		{"Non distribué", model.StatusFailed},
		{"Retour expéditeur", model.StatusFailed},
	}

	for _, tt := range tests {
		got := mapColisPriveStatus(tt.label)
		if got != tt.status {
			t.Errorf("mapColisPriveStatus(%q) = %q, want %q", tt.label, got, tt.status)
		}
	}
}

func TestParseColisPriveDate(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"01/06/2025 10:00", true},
		{"01/06/2025 10h00", true},
		{"01/06/2025", true},
		{"2025-06-01 10:00:00", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		_, err := parseColisPriveDate(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("parseColisPriveDate(%q): err=%v, wantOK=%v", tt.input, err, tt.ok)
		}
	}
}
