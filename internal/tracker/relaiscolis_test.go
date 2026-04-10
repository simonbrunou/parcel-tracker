package tracker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testRelaisColisResponse = `{
  "Colis": {
    "Colis": {
      "Enseigne": "TEST",
      "DateLivraisonEstimee": "03/06/2025 14:00",
      "ListEvenements": {
        "Evenement": [
          {
            "Libelle": "Colis pris en charge",
            "Date": "01/06/2025 10:00",
            "CodeJUS": "PCH"
          },
          {
            "Libelle": "Colis en cours d'acheminement",
            "Date": "02/06/2025 08:30",
            "CodeJUS": "ACH"
          },
          {
            "Libelle": "Colis mis en livraison",
            "Date": "03/06/2025 07:00",
            "CodeJUS": "MLV"
          },
          {
            "Libelle": "Colis livré",
            "Date": "03/06/2025 14:23",
            "CodeJUS": "LIV"
          }
        ]
      }
    }
  }
}`

func TestRelaisColisTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/x-www-form-urlencoded" {
			t.Errorf("expected form content-type, got %s", ct)
		}

		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		if bodyStr == "" {
			t.Error("expected non-empty request body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testRelaisColisResponse))
	}))
	defer server.Close()

	result, err := parseRelaisColisResponse([]byte(testRelaisColisResponse))
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
		{0, model.StatusInfoReceived, "Colis pris en charge"},
		{1, model.StatusInTransit, "Colis en cours d'acheminement"},
		{2, model.StatusOutForDelivery, "Colis mis en livraison"},
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

func TestRelaisColisCodeAndName(t *testing.T) {
	tracker := &RelaisColisTracker{}
	if tracker.Code() != model.CarrierRelaisColis {
		t.Errorf("expected code %q, got %q", model.CarrierRelaisColis, tracker.Code())
	}
	if tracker.Name() != "Relais Colis" {
		t.Errorf("expected name %q, got %q", "Relais Colis", tracker.Name())
	}
}

func TestParseRelaisColisResponseNoData(t *testing.T) {
	_, err := parseRelaisColisResponse([]byte(`{"Colis": null}`))
	if err == nil {
		t.Fatal("expected error for null colis data")
	}
}

func TestMapRelaisColisStatus(t *testing.T) {
	tests := []struct {
		code, label string
		status      model.ParcelStatus
	}{
		{"LIV", "Colis livré", model.StatusDelivered},
		{"PCH", "Colis pris en charge", model.StatusInfoReceived},
		{"ACH", "Colis en cours d'acheminement", model.StatusInTransit},
		{"MLV", "Colis mis en livraison", model.StatusOutForDelivery},
		{"", "Retour expéditeur", model.StatusFailed},
		{"", "Colis arrivé au dépôt", model.StatusInTransit},
	}

	for _, tt := range tests {
		got := mapRelaisColisStatus(tt.code, tt.label)
		if got != tt.status {
			t.Errorf("mapRelaisColisStatus(%q, %q) = %q, want %q", tt.code, tt.label, got, tt.status)
		}
	}
}

func TestParseRelaisColisDate(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"01/06/2025 10:00", true},
		{"01/06/2025 10:00:00", true},
		{"01/06/2025", true},
		{"2025-06-01T10:00:00+02:00", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		_, err := parseRelaisColisDate(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("parseRelaisColisDate(%q): err=%v, wantOK=%v", tt.input, err, tt.ok)
		}
	}
}
