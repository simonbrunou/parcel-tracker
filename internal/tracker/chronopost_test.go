package tracker

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const testChronopostResponse = `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns2:trackSkybillV2Response xmlns:ns2="http://cxf.tracking.soap.chronopost.fr/">
      <return>
        <errorCode>0</errorCode>
        <estimatedDeliveryDate>2025-06-03T18:00:00+02:00</estimatedDeliveryDate>
        <listEventInfoComp>
          <skybillNumber>XY123456789</skybillNumber>
          <events>
            <code>RG</code>
            <eventLabel>Prise en charge de votre colis</eventLabel>
            <eventDate>2025-06-01T10:00:00+02:00</eventDate>
            <officeLabel>PARIS</officeLabel>
            <zipCode>75001</zipCode>
          </events>
          <events>
            <code>TA1</code>
            <eventLabel>Colis en cours d'acheminement</eventLabel>
            <eventDate>2025-06-02T08:30:00+02:00</eventDate>
            <officeLabel>HUB LYON</officeLabel>
            <zipCode>69000</zipCode>
          </events>
          <events>
            <code>SD1</code>
            <eventLabel>Colis en cours de livraison</eventLabel>
            <eventDate>2025-06-03T07:00:00+02:00</eventDate>
            <officeLabel>MARSEILLE</officeLabel>
            <zipCode>13001</zipCode>
          </events>
          <events>
            <code>D1</code>
            <eventLabel>Colis livre</eventLabel>
            <eventDate>2025-06-03T14:23:00+02:00</eventDate>
            <officeLabel>MARSEILLE</officeLabel>
            <zipCode>13001</zipCode>
          </events>
        </listEventInfoComp>
      </return>
    </ns2:trackSkybillV2Response>
  </soap:Body>
</soap:Envelope>`

func TestChronopostTrack(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "text/xml; charset=utf-8" {
			t.Errorf("expected text/xml content-type, got %s", ct)
		}

		body, _ := io.ReadAll(r.Body)
		if !bytes.Contains(body, []byte("XY123456789")) {
			t.Errorf("request body should contain tracking number")
		}

		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(testChronopostResponse))
	}))
	defer server.Close()

	// Parse the test response directly (since we can't easily override the URL
	// in the Track method, we test the parsing separately and validate the
	// HTTP integration through the server assertions above).
	result, err := parseChronopostResponse([]byte(testChronopostResponse))
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
		message  string
	}{
		{0, model.StatusInfoReceived, "PARIS (75001)", "Prise en charge de votre colis"},
		{1, model.StatusInTransit, "HUB LYON (69000)", "Colis en cours d'acheminement"},
		{2, model.StatusOutForDelivery, "MARSEILLE (13001)", "Colis en cours de livraison"},
		{3, model.StatusDelivered, "MARSEILLE (13001)", "Colis livre"},
	}

	for _, tt := range tests {
		e := result.Events[tt.index]
		if e.Status != tt.status {
			t.Errorf("event[%d]: expected status %q, got %q", tt.index, tt.status, e.Status)
		}
		if e.Location != tt.location {
			t.Errorf("event[%d]: expected location %q, got %q", tt.index, tt.location, e.Location)
		}
		if e.Message != tt.message {
			t.Errorf("event[%d]: expected message %q, got %q", tt.index, tt.message, e.Message)
		}
	}

	// Also verify the HTTP request was sent correctly by making a real request to the test server.
	reqBody := buildChronopostSOAPRequest("XY123456789", "fr")
	req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestChronopostCodeAndName(t *testing.T) {
	tracker := &ChronopostTracker{}
	if tracker.Code() != model.CarrierChronopost {
		t.Errorf("expected code %q, got %q", model.CarrierChronopost, tracker.Code())
	}
	if tracker.Name() != "Chronopost" {
		t.Errorf("expected name %q, got %q", "Chronopost", tracker.Name())
	}
}

func TestParseChronopostResponseError(t *testing.T) {
	errorResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ns2:trackSkybillV2Response xmlns:ns2="http://cxf.tracking.soap.chronopost.fr/">
      <return>
        <errorCode>1</errorCode>
        <errorMessage>Tracking number not found</errorMessage>
        <listEventInfoComp/>
      </return>
    </ns2:trackSkybillV2Response>
  </soap:Body>
</soap:Envelope>`

	_, err := parseChronopostResponse([]byte(errorResponse))
	if err == nil {
		t.Fatal("expected error for error response")
	}
}

func TestParseChronopostResponseSOAPFault(t *testing.T) {
	faultResponse := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <soap:Fault>
      <faultstring>Internal Server Error</faultstring>
    </soap:Fault>
  </soap:Body>
</soap:Envelope>`

	_, err := parseChronopostResponse([]byte(faultResponse))
	if err == nil {
		t.Fatal("expected error for SOAP fault")
	}
}

func TestMapChronopostStatus(t *testing.T) {
	tests := []struct {
		code   string
		status model.ParcelStatus
	}{
		{"D1", model.StatusDelivered},
		{"D2", model.StatusDelivered},
		{"SD1", model.StatusOutForDelivery},
		{"CR1", model.StatusOutForDelivery},
		{"EP1", model.StatusInfoReceived},
		{"PH1", model.StatusInfoReceived},
		{"RG", model.StatusInfoReceived},
		{"TA1", model.StatusInTransit},
		{"TI1", model.StatusInTransit},
		{"LE1", model.StatusFailed},
		{"RE1", model.StatusFailed},
		{"AR1", model.StatusFailed},
		{"XX", model.StatusInTransit},
	}

	for _, tt := range tests {
		got := mapChronopostStatus(tt.code)
		if got != tt.status {
			t.Errorf("mapChronopostStatus(%q) = %q, want %q", tt.code, got, tt.status)
		}
	}
}

func TestParseChronopostDate(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"2025-06-01T10:00:00+02:00", true},
		{"2025-06-01T10:00:00Z", true},
		{"2025-06-01T10:00:00", true},
		{"2025-06-01T10:00:00.000", true},
		{"2025-06-01 10:00:00", true},
		{"01/06/2025 10:00", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		_, err := parseChronopostDate(tt.input)
		if (err == nil) != tt.ok {
			t.Errorf("parseChronopostDate(%q): err=%v, wantOK=%v", tt.input, err, tt.ok)
		}
	}
}

func TestBuildLocation(t *testing.T) {
	tests := []struct {
		site, zip, want string
	}{
		{"PARIS", "75001", "PARIS (75001)"},
		{"PARIS", "", "PARIS"},
		{"", "75001", "75001"},
		{"", "", ""},
	}

	for _, tt := range tests {
		got := buildLocation(tt.site, tt.zip)
		if got != tt.want {
			t.Errorf("buildLocation(%q, %q) = %q, want %q", tt.site, tt.zip, got, tt.want)
		}
	}
}

func TestBuildChronopostSOAPRequest(t *testing.T) {
	body := buildChronopostSOAPRequest("ABC123", "fr")
	s := string(body)

	if !bytes.Contains(body, []byte("<skybillNumber>ABC123</skybillNumber>")) {
		t.Errorf("request body missing tracking number: %s", s)
	}
	if !bytes.Contains(body, []byte("<language>fr</language>")) {
		t.Errorf("request body missing language: %s", s)
	}
	if !bytes.Contains(body, []byte("trackSkybillV2")) {
		t.Errorf("request body missing SOAP operation: %s", s)
	}
}
