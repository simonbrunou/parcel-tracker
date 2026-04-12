package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testVintedGoResponse = `{
  "tracking_code": "VTDGO123456789",
  "estimated_delivery": "2025-06-03T18:00:00Z",
  "tracking_events": [
    {
      "status": "PICKED_UP",
      "description": "Parcel picked up from sender",
      "timestamp": "2025-06-01T10:00:00Z",
      "location": {
        "city": "Vilnius",
        "country_code": "LT"
      }
    },
    {
      "status": "IN_TRANSIT",
      "description": "Parcel is in transit",
      "timestamp": "2025-06-02T08:30:00Z",
      "location": {
        "city": "Warsaw",
        "country_code": "PL"
      }
    },
    {
      "status": "OUT_FOR_DELIVERY",
      "description": "Parcel is out for delivery",
      "timestamp": "2025-06-03T07:00:00Z",
      "location": {
        "city": "Paris",
        "country_code": "FR"
      }
    },
    {
      "status": "DELIVERED",
      "description": "Parcel delivered to recipient",
      "timestamp": "2025-06-03T14:23:00Z",
      "location": {
        "city": "Paris",
        "country_code": "FR"
      }
    }
  ]
}`

func TestVintedGoTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("expected Accept application/json, got %q", accept)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(testVintedGoResponse))
	}))
	defer server.Close()

	tracker := &VintedGoTracker{
		Client: server.Client(),
	}

	result, err := parseVintedGoResponse([]byte(testVintedGoResponse))
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
		index    int
		status   model.ParcelStatus
		location string
		msg      string
	}{
		{0, model.StatusInfoReceived, "Vilnius, LT", "Parcel picked up from sender"},
		{1, model.StatusInTransit, "Warsaw, PL", "Parcel is in transit"},
		{2, model.StatusOutForDelivery, "Paris, FR", "Parcel is out for delivery"},
		{3, model.StatusDelivered, "Paris, FR", "Parcel delivered to recipient"},
	}

	for _, tt := range tests {
		e := result.Events[tt.index]
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

	_ = tracker // verify tracker was created with test client
}

func TestVintedGoCodeAndName(t *testing.T) {
	tracker := &VintedGoTracker{}
	if tracker.Code() != model.CarrierVintedGo {
		t.Errorf("expected code %q, got %q", model.CarrierVintedGo, tracker.Code())
	}
	if tracker.Name() != "Vinted Go" {
		t.Errorf("expected name %q, got %q", "Vinted Go", tracker.Name())
	}
}

func TestParseVintedGoResponseEmpty(t *testing.T) {
	result, err := parseVintedGoResponse([]byte(`{"tracking_events": []}`))
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

func TestParseVintedGoResponseInvalid(t *testing.T) {
	_, err := parseVintedGoResponse([]byte(`{invalid json}`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMapVintedGoStatus(t *testing.T) {
	tests := []struct {
		status string
		want   model.ParcelStatus
	}{
		{"DELIVERED", model.StatusDelivered},
		{"COMPLETED", model.StatusDelivered},
		{"DELIVERED_TO_PICKUP_POINT", model.StatusDelivered},
		{"OUT_FOR_DELIVERY", model.StatusOutForDelivery},
		{"DELIVERING", model.StatusOutForDelivery},
		{"LAST_MILE", model.StatusOutForDelivery},
		{"CREATED", model.StatusInfoReceived},
		{"LABEL_CREATED", model.StatusInfoReceived},
		{"REGISTERED", model.StatusInfoReceived},
		{"PICKED_UP", model.StatusInfoReceived},
		{"COLLECTED", model.StatusInfoReceived},
		{"HANDED_OVER", model.StatusInfoReceived},
		{"IN_TRANSIT", model.StatusInTransit},
		{"TRANSIT", model.StatusInTransit},
		{"SORTING", model.StatusInTransit},
		{"RETURNED", model.StatusFailed},
		{"RETURN_TO_SENDER", model.StatusFailed},
		{"FAILED", model.StatusFailed},
		{"CANCELLED", model.StatusFailed},
		{"REFUSED", model.StatusFailed},
		{"EXPIRED", model.StatusExpired},
		{"UNKNOWN_STATUS", model.StatusInTransit},
		{"", model.StatusInTransit},
	}

	for _, tt := range tests {
		got := mapVintedGoStatus(tt.status, "")
		if got != tt.want {
			t.Errorf("mapVintedGoStatus(%q) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestParseVintedGoDate(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"2025-06-01T10:00:00Z", true},
		{"2025-06-01T10:00:00+02:00", true},
		{"2025-06-01T10:00:00", true},
		{"2025-06-01T10:00:00.000", true},
		{"2025-06-01 10:00:00", true},
		{"2025-06-01", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		_, err := parseVintedGoDate(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("parseVintedGoDate(%q): err=%v, wantOK=%v", tt.input, err, tt.ok)
		}
	}
}

func TestBuildVintedGoLocation(t *testing.T) {
	tests := []struct {
		loc  *vintedGoLocation
		want string
	}{
		{&vintedGoLocation{City: "Paris", CountryCode: "FR"}, "Paris, FR"},
		{&vintedGoLocation{City: "Paris"}, "Paris"},
		{&vintedGoLocation{CountryCode: "FR"}, "FR"},
		{&vintedGoLocation{}, ""},
		{nil, ""},
	}

	for _, tt := range tests {
		got := buildVintedGoLocation(tt.loc)
		if got != tt.want {
			t.Errorf("buildVintedGoLocation(%+v) = %q, want %q", tt.loc, got, tt.want)
		}
	}
}
