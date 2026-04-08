package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testGLSResponse = `{
  "tuStatus": [
    {
      "infos": [
        {"type": "PRODUCT", "value": "Business Parcel"},
        {"type": "WEIGHT", "value": "2.5 kg"}
      ],
      "history": [
        {
          "evtDscr": "Le colis a été enregistré par l'expéditeur",
          "date": "2025-06-01",
          "time": "10:00:00",
          "address": {
            "city": "PARIS",
            "countryName": "France",
            "countryCode": "FR"
          }
        },
        {
          "evtDscr": "Le colis est en transit",
          "date": "2025-06-02",
          "time": "08:30:00",
          "address": {
            "city": "LYON",
            "countryName": "France",
            "countryCode": "FR"
          }
        },
        {
          "evtDscr": "Le colis est en cours de livraison",
          "date": "2025-06-03",
          "time": "07:00:00",
          "address": {
            "city": "MARSEILLE",
            "countryName": "France",
            "countryCode": "FR"
          }
        },
        {
          "evtDscr": "Le colis a été livré",
          "date": "2025-06-03",
          "time": "14:23:00",
          "address": {
            "city": "MARSEILLE",
            "countryName": "France",
            "countryCode": "FR"
          }
        }
      ]
    }
  ]
}`

func TestGLSTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("match") == "" {
			t.Error("expected match query parameter")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testGLSResponse))
	}))
	defer server.Close()

	events, err := parseGLSResponse([]byte(testGLSResponse))
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
	}{
		{0, model.StatusInfoReceived, "PARIS, France"},
		{1, model.StatusInTransit, "LYON, France"},
		{2, model.StatusOutForDelivery, "MARSEILLE, France"},
		{3, model.StatusDelivered, "MARSEILLE, France"},
	}

	for _, tt := range tests {
		e := events[tt.index]
		if e.Status != tt.status {
			t.Errorf("event[%d]: expected status %q, got %q", tt.index, tt.status, e.Status)
		}
		if e.Location != tt.location {
			t.Errorf("event[%d]: expected location %q, got %q", tt.index, tt.location, e.Location)
		}
	}
}

func TestGLSCodeAndName(t *testing.T) {
	tracker := &GLSTracker{}
	if tracker.Code() != model.CarrierGLS {
		t.Errorf("expected code %q, got %q", model.CarrierGLS, tracker.Code())
	}
	if tracker.Name() != "GLS" {
		t.Errorf("expected name %q, got %q", "GLS", tracker.Name())
	}
}

func TestParseGLSResponseEmpty(t *testing.T) {
	_, err := parseGLSResponse([]byte(`{"tuStatus": []}`))
	if err == nil {
		t.Fatal("expected error for empty tuStatus")
	}
}

func TestMapGLSStatus(t *testing.T) {
	tests := []struct {
		desc   string
		status model.ParcelStatus
	}{
		{"Le colis a été livré", model.StatusDelivered},
		{"Delivered", model.StatusDelivered},
		{"Le colis est en cours de livraison", model.StatusOutForDelivery},
		{"Le colis a été pris en charge", model.StatusInfoReceived},
		{"Le colis est en transit", model.StatusInTransit},
		{"Non distribuable", model.StatusFailed},
		{"Retour à l'expéditeur", model.StatusFailed},
	}

	for _, tt := range tests {
		got := mapGLSStatus(tt.desc)
		if got != tt.status {
			t.Errorf("mapGLSStatus(%q) = %q, want %q", tt.desc, got, tt.status)
		}
	}
}

func TestParseGLSDate(t *testing.T) {
	tests := []struct {
		date, time string
		ok         bool
	}{
		{"2025-06-01", "10:00:00", true},
		{"2025-06-01", "10:00", true},
		{"01/06/2025", "10:00:00", true},
		{"invalid", "10:00", false},
	}

	for _, tt := range tests {
		_, err := parseGLSDate(tt.date, tt.time)
		if (err == nil) != tt.ok {
			t.Errorf("parseGLSDate(%q, %q): err=%v, wantOK=%v", tt.date, tt.time, err, tt.ok)
		}
	}
}

func TestBuildGLSLocation(t *testing.T) {
	tests := []struct {
		addr glsAddress
		want string
	}{
		{glsAddress{City: "PARIS", CountryName: "France"}, "PARIS, France"},
		{glsAddress{City: "PARIS"}, "PARIS"},
		{glsAddress{CountryName: "France"}, "France"},
		{glsAddress{}, ""},
	}

	for _, tt := range tests {
		got := buildGLSLocation(tt.addr)
		if got != tt.want {
			t.Errorf("buildGLSLocation(%+v) = %q, want %q", tt.addr, got, tt.want)
		}
	}
}
