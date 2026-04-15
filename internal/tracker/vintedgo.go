package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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

// Vinted Go JSON response structures — matching the actual v3 public API.

type vintedGoResponse struct {
	TrackingEvents []vintedGoEvent `json:"tracking_events"`
	Meta           *vintedGoMeta   `json:"meta"`
}

type vintedGoEvent struct {
	ID            int                `json:"id"`
	Message       string             `json:"message"`
	Timestamp     string             `json:"timestamp"`
	Metadata      vintedGoMetadata   `json:"metadata"`
	BannerMessage *string            `json:"banner_message"`
	GroupHeader   *string            `json:"group_header"`
	Group         string             `json:"group"`
	State         string             `json:"state"`
}

type vintedGoMetadata struct {
	Address string `json:"address,omitempty"`
	PointID int    `json:"point_id,omitempty"`
}

type vintedGoMeta struct {
	Banner         *string `json:"banner"`
	ExpirationTime string  `json:"expiration_time"`
}

func parseVintedGoResponse(data []byte) (TrackResult, error) {
	var resp vintedGoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return TrackResult{}, fmt.Errorf("vintedgo: parse json: %w", err)
	}

	var result TrackResult

	for _, e := range resp.TrackingEvents {
		ts, err := parseVintedGoDate(e.Timestamp)
		if err != nil {
			continue
		}

		result.Events = append(result.Events, model.TrackingEvent{
			Status:    mapVintedGoStatus(e.Group, e.State),
			Message:   e.Message,
			Location:  extractVintedGoLocation(e),
			Timestamp: ts.UTC(),
		})
	}

	return result, nil
}

func parseVintedGoDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
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

// locationInPattern matches messages like "At sorting center in Paris, FR".
var locationInPattern = regexp.MustCompile(`\bin\s+(.+)$`)

// extractVintedGoLocation derives a location string from event metadata or message.
func extractVintedGoLocation(e vintedGoEvent) string {
	// Prefer explicit address from metadata.
	if addr := strings.TrimSpace(e.Metadata.Address); addr != "" {
		return addr
	}

	// Fall back to parsing location from the message text (e.g. "At sorting center in Paris, FR").
	if m := locationInPattern.FindStringSubmatch(e.Message); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}

	return ""
}

// mapVintedGoStatus maps a Vinted Go event group/state to an internal ParcelStatus.
// The "group" field is the primary discriminator in the v3 API.
func mapVintedGoStatus(group, state string) model.ParcelStatus {
	switch strings.ToLower(strings.TrimSpace(group)) {
	case "created":
		return model.StatusInfoReceived
	case "shipped", "handed_over", "collected":
		return model.StatusInfoReceived
	case "in_transit", "transit", "sorting":
		return model.StatusInTransit
	case "out_for_delivery", "delivering", "last_mile":
		return model.StatusOutForDelivery
	case "ready_for_pickup", "ready_for_collection":
		return model.StatusOutForDelivery
	case "delivered", "completed":
		return model.StatusDelivered
	case "returned", "return_to_sender", "failed", "cancelled", "refused":
		return model.StatusFailed
	case "expired":
		return model.StatusExpired
	default:
		return model.StatusInTransit
	}
}
