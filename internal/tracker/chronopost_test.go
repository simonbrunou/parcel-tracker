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
            <code>LT</code>
            <eventLabel>Colis en cours de livraison</eventLabel>
            <eventDate>2025-06-03T07:00:00+02:00</eventDate>
            <officeLabel>MARSEILLE</officeLabel>
            <zipCode>13001</zipCode>
          </events>
          <events>
            <code>LV</code>
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
		label  string
		status model.ParcelStatus
	}{
		// Delivered (by code)
		{"LV", "", model.StatusDelivered},
		{"RM", "", model.StatusDelivered},
		{"BL", "", model.StatusDelivered},
		{"LP", "", model.StatusDelivered},
		{"D1", "", model.StatusDelivered},
		{"D2", "", model.StatusDelivered},
		// Out for delivery
		{"LT", "", model.StatusOutForDelivery},
		{"CR1", "", model.StatusOutForDelivery},
		{"MD1", "", model.StatusOutForDelivery},
		// Info received / preparation / pickup
		{"DC", "", model.StatusInfoReceived},
		{"EP1", "", model.StatusInfoReceived},
		{"PH1", "", model.StatusInfoReceived},
		{"RG", "", model.StatusInfoReceived},
		// In transit
		{"EC", "", model.StatusInTransit},
		{"TS", "", model.StatusInTransit},
		{"SD", "", model.StatusInTransit}, // sorted at depot, NOT out for delivery
		{"IS", "", model.StatusInTransit},
		{"TA1", "", model.StatusInTransit},
		{"TI1", "", model.StatusInTransit},
		// Padded with whitespace (Chronopost real-world format)
		{"DC ", "", model.StatusInfoReceived},
		{"SD ", "", model.StatusInTransit},
		// Failed
		{"LE1", "", model.StatusFailed},
		{"RE1", "", model.StatusFailed},
		{"AR1", "", model.StatusFailed},
		// Unknown code, no label
		{"XX", "", model.StatusInTransit},

		// Label-based fallback for unknown codes
		{"XX", "Colis livré", model.StatusDelivered},
		{"XX", "Remis au destinataire", model.StatusDelivered},
		{"XX", "Distribué", model.StatusDelivered},
		{"XX", "Package delivered", model.StatusDelivered},
		{"XX", "Colis en cours de livraison", model.StatusOutForDelivery},
		{"XX", "Mis en livraison", model.StatusOutForDelivery},
		{"XX", "Retour à l'expéditeur", model.StatusFailed},
		{"XX", "Refusé par le destinataire", model.StatusFailed},
		{"XX", "Anomalie constatée", model.StatusFailed},
		{"XX", "Non distribué", model.StatusFailed},
		{"XX", "Colis en cours d'acheminement", model.StatusInTransit},

		// Code takes priority over label: DC is InfoReceived even if label contains "livr"
		{"DC", "Colis livré au voisin", model.StatusInfoReceived},
	}

	for _, tt := range tests {
		got := mapChronopostStatus(tt.code, tt.label)
		if got != tt.status {
			t.Errorf("mapChronopostStatus(%q, %q) = %q, want %q", tt.code, tt.label, got, tt.status)
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

// testChronopostLiveResponse is an actual response captured from the live
// Chronopost tracking-cxf API for tracking number XG953488284JB. It exercises
// real-world quirks: codes padded with trailing spaces, whitespace-only zip
// codes, the DC preparation code (which must not be classified as Delivered),
// and the SD sorted-at-depot code (which must not be classified as
// OutForDelivery).
const testChronopostLiveResponse = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><ns2:trackSkybillV2Response xmlns:ns2="http://cxf.tracking.soap.chronopost.fr/"><return><errorCode>0</errorCode><listEventInfoComp><events><code>DC </code><eventDate>2026-04-10T07:55:27+02:00</eventDate><eventLabel>Colis en cours de pr&#233;paration chez l'exp&#233;diteur</eventLabel><highPriority>false</highPriority><NPC>   </NPC><officeLabel>Web Services</officeLabel><zipCode>     </zipCode></events><events><code>EC </code><eventDate>2026-04-10T14:28:01+02:00</eventDate><eventLabel>Tri effectu&#233; dans l'agence de d&#233;part</eventLabel><highPriority>true</highPriority><NPC> 62</NPC><officeLabel>ARRAS CHRONOPOST</officeLabel><zipCode>62223</zipCode></events><events><code>TS </code><eventDate>2026-04-10T14:29:01+02:00</eventDate><eventLabel>Colis en cours d'acheminement</eventLabel><highPriority>false</highPriority><NPC> 62</NPC><officeLabel>ARRAS CHRONOPOST</officeLabel><zipCode>62223</zipCode></events><events><code>TS </code><eventDate>2026-04-10T19:18:13+02:00</eventDate><eventLabel>Colis en cours d'acheminement</eventLabel><highPriority>false</highPriority><NPC>93 </NPC><officeLabel>AULNAY SOUS BOIS CHRONOPOST</officeLabel><zipCode>93600</zipCode></events><events><code>SD </code><eventDate>2026-04-11T05:26:24+02:00</eventDate><eventLabel>Tri effectu&#233; dans l'agence de distribution</eventLabel><highPriority>false</highPriority><NPC>56 </NPC><officeLabel>VANNES CHRONOPOST</officeLabel><zipCode>56000</zipCode></events><events><code>IS </code><eventDate>2026-04-11T05:27:24+02:00</eventDate><eventLabel>Livraison pr&#233;vue lundi prochain</eventLabel><highPriority>false</highPriority><NPC>56 </NPC><officeLabel>VANNES CHRONOPOST</officeLabel><zipCode>56000</zipCode></events><skybillNumber>XG953488284JB</skybillNumber></listEventInfoComp></return></ns2:trackSkybillV2Response></soap:Body></soap:Envelope>`

// TestParseChronopostLiveResponse pins the parser against a real Chronopost
// response so status regressions are caught. In particular it guards the
// DC -> InfoReceived and SD -> InTransit mappings, which were previously
// misclassified as Delivered and OutForDelivery respectively.
func TestParseChronopostLiveResponse(t *testing.T) {
	result, err := parseChronopostResponse([]byte(testChronopostLiveResponse))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Events) != 6 {
		t.Fatalf("expected 6 events, got %d", len(result.Events))
	}

	// No estimatedDeliveryDate in this response.
	if result.EstimatedDelivery != nil {
		t.Errorf("expected no estimated delivery, got %v", result.EstimatedDelivery)
	}

	want := []struct {
		status   model.ParcelStatus
		location string
	}{
		{model.StatusInfoReceived, "Web Services"},            // DC
		{model.StatusInTransit, "ARRAS CHRONOPOST (62223)"},    // EC
		{model.StatusInTransit, "ARRAS CHRONOPOST (62223)"},    // TS
		{model.StatusInTransit, "AULNAY SOUS BOIS CHRONOPOST (93600)"}, // TS
		{model.StatusInTransit, "VANNES CHRONOPOST (56000)"},   // SD (NOT out for delivery)
		{model.StatusInTransit, "VANNES CHRONOPOST (56000)"},   // IS
	}

	for i, w := range want {
		e := result.Events[i]
		if e.Status != w.status {
			t.Errorf("event[%d]: expected status %q, got %q", i, w.status, e.Status)
		}
		if e.Location != w.location {
			t.Errorf("event[%d]: expected location %q, got %q", i, w.location, e.Location)
		}
		if e.Timestamp.IsZero() {
			t.Errorf("event[%d]: timestamp should be parsed", i)
		}
	}

	// Guard against the previous bug where DC was classified as Delivered
	// because HasPrefix("D") greedy-matched it.
	if result.Events[0].Status == model.StatusDelivered {
		t.Errorf("DC event must not be classified as Delivered")
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
