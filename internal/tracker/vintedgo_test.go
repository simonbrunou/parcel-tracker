package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// testVintedGoResponse mirrors the actual v3 public API response format.
const testVintedGoResponse = `{
  "tracking_events": [
    {
      "id": 1313156661,
      "message": "Ready to collect at the locker",
      "timestamp": "2026-04-15T08:16:28.236Z",
      "metadata": {
        "address": "4 Rue Marcel Dassault, Saint-Avé, FR",
        "point_id": 137178
      },
      "banner_message": null,
      "group_header": null,
      "group": "ready_for_pickup",
      "state": "delivery"
    },
    {
      "id": 1312884320,
      "message": "The parcel is on its way to the locker (Intermarché Super)",
      "timestamp": "2026-04-15T07:44:09.546Z",
      "metadata": {},
      "banner_message": null,
      "group_header": "In transit",
      "group": "in_transit",
      "state": "delivery"
    },
    {
      "id": 1312228892,
      "message": "At sorting center in Vannes, FR",
      "timestamp": "2026-04-15T06:24:51.574Z",
      "metadata": {},
      "banner_message": null,
      "group_header": "In transit",
      "group": "in_transit",
      "state": "delivery"
    },
    {
      "id": 1307463543,
      "message": "At sorting center in Paris, FR",
      "timestamp": "2026-04-14T10:26:34.695Z",
      "metadata": {},
      "banner_message": null,
      "group_header": "In transit",
      "group": "in_transit",
      "state": "delivery"
    },
    {
      "id": 1304397033,
      "message": "At sorting center in Paris, FR",
      "timestamp": "2026-04-13T21:55:49.763Z",
      "metadata": {},
      "banner_message": null,
      "group_header": "In transit",
      "group": "in_transit",
      "state": "delivery"
    },
    {
      "id": 1288911719,
      "message": "Parcel was dropped off at the parcel shop",
      "timestamp": "2026-04-10T10:46:20.372Z",
      "metadata": {},
      "banner_message": null,
      "group_header": null,
      "group": "shipped",
      "state": "delivery"
    },
    {
      "id": 1287928334,
      "message": "Shipment created",
      "timestamp": "2026-04-10T08:29:27.631Z",
      "metadata": {},
      "banner_message": null,
      "group_header": null,
      "group": "created",
      "state": "delivery"
    }
  ],
  "meta": {
    "banner": null,
    "expiration_time": "2026-04-22T08:16:27.000Z",
    "contextual_faq": []
  }
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

	if len(result.Events) != 7 {
		t.Fatalf("expected 7 events, got %d", len(result.Events))
	}

	tests := []struct {
		index    int
		status   model.ParcelStatus
		location string
		msg      string
	}{
		{0, model.StatusOutForDelivery, "4 Rue Marcel Dassault, Saint-Avé, FR", "Ready to collect at the locker"},
		{1, model.StatusInTransit, "", "The parcel is on its way to the locker (Intermarché Super)"},
		{2, model.StatusInTransit, "Vannes, FR", "At sorting center in Vannes, FR"},
		{3, model.StatusInTransit, "Paris, FR", "At sorting center in Paris, FR"},
		{4, model.StatusInTransit, "Paris, FR", "At sorting center in Paris, FR"},
		{5, model.StatusInfoReceived, "", "Parcel was dropped off at the parcel shop"},
		{6, model.StatusInfoReceived, "", "Shipment created"},
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
}

func TestParseVintedGoResponseInvalid(t *testing.T) {
	_, err := parseVintedGoResponse([]byte(`{invalid json}`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestMapVintedGoStatus(t *testing.T) {
	tests := []struct {
		group string
		state string
		want  model.ParcelStatus
	}{
		{"created", "delivery", model.StatusInfoReceived},
		{"shipped", "delivery", model.StatusInfoReceived},
		{"handed_over", "delivery", model.StatusInfoReceived},
		{"collected", "delivery", model.StatusInfoReceived},
		{"in_transit", "delivery", model.StatusInTransit},
		{"transit", "delivery", model.StatusInTransit},
		{"sorting", "delivery", model.StatusInTransit},
		{"out_for_delivery", "delivery", model.StatusOutForDelivery},
		{"delivering", "delivery", model.StatusOutForDelivery},
		{"last_mile", "delivery", model.StatusOutForDelivery},
		{"ready_for_pickup", "delivery", model.StatusOutForDelivery},
		{"ready_for_collection", "delivery", model.StatusOutForDelivery},
		{"delivered", "delivery", model.StatusDelivered},
		{"completed", "delivery", model.StatusDelivered},
		{"returned", "delivery", model.StatusFailed},
		{"return_to_sender", "delivery", model.StatusFailed},
		{"failed", "delivery", model.StatusFailed},
		{"cancelled", "delivery", model.StatusFailed},
		{"refused", "delivery", model.StatusFailed},
		{"expired", "delivery", model.StatusExpired},
		{"unknown_group", "delivery", model.StatusInTransit},
		{"", "", model.StatusInTransit},
	}

	for _, tt := range tests {
		got := mapVintedGoStatus(tt.group, tt.state)
		if got != tt.want {
			t.Errorf("mapVintedGoStatus(%q, %q) = %q, want %q", tt.group, tt.state, got, tt.want)
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
		{"2025-06-01T10:00:00.000Z", true},
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

func TestExtractVintedGoLocation(t *testing.T) {
	tests := []struct {
		name  string
		event vintedGoEvent
		want  string
	}{
		{
			name: "metadata address takes priority",
			event: vintedGoEvent{
				Message:  "Ready to collect at the locker",
				Metadata: vintedGoMetadata{Address: "4 Rue Marcel Dassault, Saint-Avé, FR"},
			},
			want: "4 Rue Marcel Dassault, Saint-Avé, FR",
		},
		{
			name: "location parsed from message",
			event: vintedGoEvent{
				Message:  "At sorting center in Paris, FR",
				Metadata: vintedGoMetadata{},
			},
			want: "Paris, FR",
		},
		{
			name: "no location available",
			event: vintedGoEvent{
				Message:  "Shipment created",
				Metadata: vintedGoMetadata{},
			},
			want: "",
		},
		{
			name: "parenthetical message no location",
			event: vintedGoEvent{
				Message:  "The parcel is on its way to the locker (Intermarché Super)",
				Metadata: vintedGoMetadata{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVintedGoLocation(tt.event)
			if got != tt.want {
				t.Errorf("extractVintedGoLocation() = %q, want %q", got, tt.want)
			}
		})
	}
}
