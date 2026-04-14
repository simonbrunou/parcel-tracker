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

const vintedGoTrackingURL = "https://carrier.vintedgo.com/public/v3/tracking_events/"

// VintedGoTracker tracks parcels via the Vinted Go public tracking API.
type VintedGoTracker struct {
	Client *http.Client
}

func (t *VintedGoTracker) Code() model.CarrierCode { return model.CarrierVintedGo }
func (t *VintedGoTracker) Name() string             { return "Vinted Go" }

func (t *VintedGoTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *VintedGoTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	url := vintedGoTrackingURL + trackingNumber

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResult{}, fmt.Errorf("vintedgo: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("vintedgo: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return TrackResult{}, fmt.Errorf("vintedgo: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("vintedgo: unexpected status %d", resp.StatusCode)
	}

	return parseVintedGoResponse(body)
}

// Vinted Go JSON response structures.

type vintedGoResponse struct {
	TrackingCode      string          `json:"tracking_code"`
	EstimatedDelivery string          `json:"estimated_delivery"`
	TrackingEvents    []vintedGoEvent `json:"tracking_events"`
}

type vintedGoEvent struct {
	Status      string            `json:"status"`
	Description string            `json:"description"`
	Timestamp   string            `json:"timestamp"`
	Location    *vintedGoLocation `json:"location"`
}

type vintedGoLocation struct {
	City        string `json:"city"`
	CountryCode string `json:"country_code"`
}

func parseVintedGoResponse(data []byte) (TrackResult, error) {
	var resp vintedGoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return TrackResult{}, fmt.Errorf("vintedgo: parse json: %w", err)
	}

	var result TrackResult

	if resp.EstimatedDelivery != "" {
		if t, err := parseVintedGoDate(resp.EstimatedDelivery); err == nil {
			utc := t.UTC()
			result.EstimatedDelivery = &utc
		}
	}

	for _, e := range resp.TrackingEvents {
		ts, err := parseVintedGoDate(e.Timestamp)
		if err != nil {
			continue
		}

		result.Events = append(result.Events, model.TrackingEvent{
			Status:    mapVintedGoStatus(e.Status, e.Description),
			Message:   e.Description,
			Location:  buildVintedGoLocation(e.Location),
			Timestamp: ts.UTC(),
		})
	}

	return result, nil
}

func parseVintedGoDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	s = strings.TrimSpace(s)
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("vintedgo: unknown date format: %q", s)
}

func buildVintedGoLocation(loc *vintedGoLocation) string {
	if loc == nil {
		return ""
	}
	city := strings.TrimSpace(loc.City)
	country := strings.TrimSpace(loc.CountryCode)
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

// mapVintedGoStatus maps a Vinted Go event status to an internal ParcelStatus.
// Uses contains-matching on the uppercase status string for resilience against
// variations like "DELIVERED_TO_PICKUP_POINT".
func mapVintedGoStatus(status, description string) model.ParcelStatus {
	upper := strings.ToUpper(strings.TrimSpace(status))

	switch {
	case strings.Contains(upper, "OUT_FOR_DELIVERY") ||
		strings.Contains(upper, "DELIVERING") ||
		strings.Contains(upper, "LAST_MILE"):
		return model.StatusOutForDelivery
	case strings.Contains(upper, "DELIVER") ||
		strings.Contains(upper, "COMPLETED"):
		return model.StatusDelivered
	case strings.Contains(upper, "CREATED") ||
		strings.Contains(upper, "REGISTERED") ||
		strings.Contains(upper, "PICKED_UP") ||
		strings.Contains(upper, "COLLECTED") ||
		strings.Contains(upper, "HANDED"):
		return model.StatusInfoReceived
	case strings.Contains(upper, "RETURN") ||
		strings.Contains(upper, "FAILED") ||
		strings.Contains(upper, "CANCEL") ||
		strings.Contains(upper, "REFUSED"):
		return model.StatusFailed
	case strings.Contains(upper, "EXPIRED"):
		return model.StatusExpired
	default:
		return model.StatusInTransit
	}
}
