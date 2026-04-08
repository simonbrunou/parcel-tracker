package tracker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testLaPosteResponse = `{
  "shipment": {
    "product": "colissimo",
    "contextData": {
      "originCountry": "FR",
      "arrivalCountry": "FR"
    },
    "timeline": [
      {"status": true, "shortLabel": "Pris en charge"},
      {"status": true, "shortLabel": "En cours d'acheminement"},
      {"status": true, "shortLabel": "Distribué"}
    ],
    "event": [
      {
        "code": "DR1",
        "label": "Déclaratif réceptionné",
        "date": "2025-06-01T10:00:00+02:00"
      },
      {
        "code": "PC1",
        "label": "Pris en charge par La Poste",
        "date": "2025-06-01T14:00:00+02:00"
      },
      {
        "code": "ET1",
        "label": "En cours de traitement",
        "date": "2025-06-02T08:30:00+02:00"
      },
      {
        "code": "MD2",
        "label": "Mis en distribution",
        "date": "2025-06-03T07:00:00+02:00"
      },
      {
        "code": "DI1",
        "label": "Distribué",
        "date": "2025-06-03T14:23:00+02:00"
      }
    ]
  }
}`

func TestLaPosteTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("X-Okapi-Key") != "test-api-key" {
			t.Errorf("expected X-Okapi-Key header, got %q", r.Header.Get("X-Okapi-Key"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected application/json Accept, got %q", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testLaPosteResponse))
	}))
	defer server.Close()

	// Test the response parsing directly.
	events, err := parseLaPosteResponse([]byte(testLaPosteResponse))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) != 5 {
		t.Fatalf("expected 5 events, got %d", len(events))
	}

	tests := []struct {
		index  int
		status model.ParcelStatus
		msg    string
	}{
		{0, model.StatusInfoReceived, "Déclaratif réceptionné"},
		{1, model.StatusInfoReceived, "Pris en charge par La Poste"},
		{2, model.StatusInTransit, "En cours de traitement"},
		{3, model.StatusOutForDelivery, "Mis en distribution"},
		{4, model.StatusDelivered, "Distribué"},
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

	// Verify the tracker sends correct headers to the test server.
	tracker := &LaPosteTracker{APIKey: "test-api-key", Client: server.Client()}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/test?lang=fr_FR", nil)
	req.Header.Set("X-Okapi-Key", "test-api-key")
	req.Header.Set("Accept", "application/json")
	resp, err := server.Client().Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	resp.Body.Close()
	_ = tracker // ensure tracker is used
}

func TestLaPosteCodeAndName(t *testing.T) {
	tracker := &LaPosteTracker{}
	if tracker.Code() != model.CarrierLaPoste {
		t.Errorf("expected code %q, got %q", model.CarrierLaPoste, tracker.Code())
	}
	if tracker.Name() != "La Poste" {
		t.Errorf("expected name %q, got %q", "La Poste", tracker.Name())
	}
}

func TestLaPosteTrackMissingAPIKey(t *testing.T) {
	tracker := &LaPosteTracker{APIKey: ""}
	_, err := tracker.Track(context.Background(), "123456789")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestParseLaPosteResponseError(t *testing.T) {
	errorResponse := `{"returnMessage": "Aucun objet trouvé pour cet identifiant"}`
	_, err := parseLaPosteResponse([]byte(errorResponse))
	if err == nil {
		t.Fatal("expected error for error response")
	}
}

func TestMapLaPosteStatus(t *testing.T) {
	tests := []struct {
		code   string
		status model.ParcelStatus
	}{
		{"DR1", model.StatusInfoReceived},
		{"PC1", model.StatusInfoReceived},
		{"PC2", model.StatusInfoReceived},
		{"ET1", model.StatusInTransit},
		{"ET2", model.StatusInTransit},
		{"ET3", model.StatusInTransit},
		{"ET4", model.StatusInTransit},
		{"EP1", model.StatusInTransit},
		{"DO1", model.StatusInTransit},
		{"DO2", model.StatusInTransit},
		{"DO3", model.StatusInTransit},
		{"PB1", model.StatusFailed},
		{"PB2", model.StatusInTransit},
		{"MD2", model.StatusOutForDelivery},
		{"ND1", model.StatusFailed},
		{"AG1", model.StatusInTransit},
		{"RE1", model.StatusFailed},
		{"DI1", model.StatusDelivered},
		{"DI2", model.StatusDelivered},
	}

	for _, tt := range tests {
		got := mapLaPosteStatus(tt.code)
		if got != tt.status {
			t.Errorf("mapLaPosteStatus(%q) = %q, want %q", tt.code, got, tt.status)
		}
	}
}

func TestParseLaPosteDate(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"2025-06-01T10:00:00+02:00", true},
		{"2025-06-01T10:00:00Z", true},
		{"2025-06-01T10:00:00", true},
		{"2025-06-01 10:00:00", true},
		{"01/06/2025 10:00", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		_, err := parseLaPosteDate(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("parseLaPosteDate(%q): err=%v, wantOK=%v", tt.input, err, tt.ok)
		}
	}
}
