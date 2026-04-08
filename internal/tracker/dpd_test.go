package tracker

import (
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testDPDHTML = `<!DOCTYPE html>
<html>
<body>
<table>
  <tr id="ligneTableTrace_0">
    <td>01/06/2025</td>
    <td>10:00</td>
    <td>Données de colis transmises</td>
    <td>PARIS</td>
  </tr>
  <tr id="ligneTableTrace_1">
    <td>02/06/2025</td>
    <td>08:30</td>
    <td>En cours d'acheminement</td>
    <td>LYON</td>
  </tr>
  <tr id="ligneTableTrace_2">
    <td>03/06/2025</td>
    <td>07:00</td>
    <td>En cours de livraison</td>
    <td>MARSEILLE</td>
  </tr>
  <tr id="ligneTableTrace_3">
    <td>03/06/2025</td>
    <td>14:23</td>
    <td>Colis livré</td>
    <td>MARSEILLE</td>
  </tr>
</table>
</body>
</html>`

func TestParseDPDHTML(t *testing.T) {
	events, err := parseDPDHTML([]byte(testDPDHTML))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}

	tests := []struct {
		index    int
		status   model.ParcelStatus
		location string
		msg      string
	}{
		{0, model.StatusInfoReceived, "PARIS", "Données de colis transmises"},
		{1, model.StatusInTransit, "LYON", "En cours d'acheminement"},
		{2, model.StatusOutForDelivery, "MARSEILLE", "En cours de livraison"},
		{3, model.StatusDelivered, "MARSEILLE", "Colis livré"},
	}

	for _, tt := range tests {
		e := events[tt.index]
		if e.Status != tt.status {
			t.Errorf("event[%d]: expected status %q, got %q", tt.index, tt.status, e.Status)
		}
		if e.Location != tt.location {
			t.Errorf("event[%d]: expected location %q, got %q", tt.index, tt.location, e.Location)
		}
		if e.Message != tt.msg {
			t.Errorf("event[%d]: expected message %q, got %q", tt.index, tt.msg, e.Message)
		}
	}
}

func TestParseDPDHTMLEmpty(t *testing.T) {
	html := `<!DOCTYPE html><html><body><p>No tracking data</p></body></html>`
	events, err := parseDPDHTML([]byte(html))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestDPDCodeAndName(t *testing.T) {
	tracker := &DPDTracker{}
	if tracker.Code() != model.CarrierDPD {
		t.Errorf("expected code %q, got %q", model.CarrierDPD, tracker.Code())
	}
	if tracker.Name() != "DPD" {
		t.Errorf("expected name %q, got %q", "DPD", tracker.Name())
	}
}

func TestMapDPDStatus(t *testing.T) {
	tests := []struct {
		label  string
		status model.ParcelStatus
	}{
		{"Colis livré", model.StatusDelivered},
		{"Remis au destinataire", model.StatusDelivered},
		{"En cours de livraison", model.StatusOutForDelivery},
		{"Avec le chauffeur pour livraison", model.StatusOutForDelivery},
		{"Données de colis transmises", model.StatusInfoReceived},
		{"Pris en charge par DPD", model.StatusInfoReceived},
		{"En cours d'acheminement", model.StatusInTransit},
		{"Non distribué - Absence", model.StatusFailed},
		{"Retour expéditeur", model.StatusFailed},
	}

	for _, tt := range tests {
		got := mapDPDStatus(tt.label)
		if got != tt.status {
			t.Errorf("mapDPDStatus(%q) = %q, want %q", tt.label, got, tt.status)
		}
	}
}

func TestCleanDPDLabel(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Predict vous informe : \nLivraison prévue", "Livraison prévue"},
		{"Instruction :faire suivre", "faire suivre"},
		{"Normal label", "Normal label"},
	}

	for _, tt := range tests {
		got := cleanDPDLabel(tt.input)
		if got != tt.want {
			t.Errorf("cleanDPDLabel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseDPDDate(t *testing.T) {
	tests := []struct {
		date, time string
		ok         bool
	}{
		{"01/06/2025", "10:00", true},
		{"01/06/2025", "10h00", true},
		{"2025-06-01", "10:00:00", true},
		{"invalid", "10:00", false},
	}

	for _, tt := range tests {
		_, err := parseDPDDate(tt.date, tt.time)
		if (err == nil) != tt.ok {
			t.Errorf("parseDPDDate(%q, %q): err=%v, wantOK=%v", tt.date, tt.time, err, tt.ok)
		}
	}
}
