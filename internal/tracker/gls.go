package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const glsTrackingURL = "https://gls-group.eu/app/service/open/rest/FR/fr/rstt001"

// GLSTracker tracks parcels via the GLS Group public REST API.
type GLSTracker struct {
	Client *http.Client
}

func (t *GLSTracker) Code() model.CarrierCode { return model.CarrierGLS }
func (t *GLSTracker) Name() string             { return "GLS" }

func (t *GLSTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *GLSTracker) Track(ctx context.Context, trackingNumber string) ([]model.TrackingEvent, error) {
	url := glsTrackingURL + "?match=" + trackingNumber

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gls: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("gls: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gls: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gls: unexpected status %d", resp.StatusCode)
	}

	return parseGLSResponse(body)
}

// GLS JSON response structures.

type glsResponse struct {
	TuStatus []glsShipment `json:"tuStatus"`
}

type glsShipment struct {
	History []glsEvent `json:"history"`
}

type glsEvent struct {
	EvtDscr string     `json:"evtDscr"`
	Date    string     `json:"date"`
	Time    string     `json:"time"`
	Address glsAddress `json:"address"`
}

type glsAddress struct {
	City        string `json:"city"`
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
}

func parseGLSResponse(data []byte) ([]model.TrackingEvent, error) {
	var resp glsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("gls: parse json: %w", err)
	}

	if len(resp.TuStatus) == 0 {
		return nil, fmt.Errorf("gls: no shipment data")
	}

	shipment := resp.TuStatus[0]
	var events []model.TrackingEvent

	for _, e := range shipment.History {
		ts, err := parseGLSDate(e.Date, e.Time)
		if err != nil {
			continue
		}

		location := buildGLSLocation(e.Address)

		events = append(events, model.TrackingEvent{
			Status:    mapGLSStatus(e.EvtDscr),
			Message:   e.EvtDscr,
			Location:  location,
			Timestamp: ts.UTC(),
		})
	}

	return events, nil
}

func parseGLSDate(date, timeStr string) (time.Time, error) {
	combined := strings.TrimSpace(date) + " " + strings.TrimSpace(timeStr)
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"02/01/2006 15:04:05",
		"02/01/2006 15:04",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, combined); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("gls: unknown date format: %q", combined)
}

func buildGLSLocation(addr glsAddress) string {
	city := strings.TrimSpace(addr.City)
	country := strings.TrimSpace(addr.CountryName)
	switch {
	case city != "" && country != "":
		return city + ", " + country
	case city != "":
		return city
	case country != "":
		return country
	default:
		return ""
	}
}

// mapGLSStatus maps a GLS event description to an internal ParcelStatus.
// GLS does not use structured status codes in its public API, so we match
// on keywords in the event description (French and English).
func mapGLSStatus(description string) model.ParcelStatus {
	lower := strings.ToLower(description)

	switch {
	// Out for delivery must be checked before delivered (both contain "livr").
	case strings.Contains(lower, "en cours de livraison") ||
		strings.Contains(lower, "out for delivery") ||
		strings.Contains(lower, "en livraison") ||
		strings.Contains(lower, "dernier kilomètre") ||
		strings.Contains(lower, "chauffeur"):
		return model.StatusOutForDelivery
	case strings.Contains(lower, "livr") ||
		strings.Contains(lower, "delivered") ||
		strings.Contains(lower, "remis"):
		return model.StatusDelivered
	case strings.Contains(lower, "pris en charge") ||
		strings.Contains(lower, "enregistr") ||
		strings.Contains(lower, "pickup") ||
		strings.Contains(lower, "collected"):
		return model.StatusInfoReceived
	case strings.Contains(lower, "non distribuable") ||
		strings.Contains(lower, "retour") ||
		strings.Contains(lower, "refus") ||
		strings.Contains(lower, "not delivered") ||
		strings.Contains(lower, "returned"):
		return model.StatusFailed
	default:
		return model.StatusInTransit
	}
}
