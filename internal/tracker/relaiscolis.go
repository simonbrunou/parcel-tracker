package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const relaisColisTrackingURL = "https://www.relaiscolis.com/suivi-de-colis/index/tracking/"

// RelaisColisTracker tracks parcels via the Relais Colis public tracking API.
type RelaisColisTracker struct {
	Client *http.Client
}

func (t *RelaisColisTracker) Code() model.CarrierCode { return model.CarrierRelaisColis }
func (t *RelaisColisTracker) Name() string             { return "Relais Colis" }

func (t *RelaisColisTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *RelaisColisTracker) Track(ctx context.Context, trackingNumber string) ([]model.TrackingEvent, error) {
	form := url.Values{
		"valeur":        {trackingNumber},
		"typeRecherche": {"EXP"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, relaisColisTrackingURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("relaiscolis: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ParcelTracker/1.0)")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("relaiscolis: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("relaiscolis: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("relaiscolis: unexpected status %d", resp.StatusCode)
	}

	return parseRelaisColisResponse(body)
}

// Relais Colis JSON response structures.

type relaisColisResponse struct {
	Colis *relaisColisWrapper `json:"Colis"`
}

type relaisColisWrapper struct {
	Colis *relaisColisData `json:"Colis"`
}

type relaisColisData struct {
	Enseigne        string                   `json:"Enseigne"`
	ListEvenements  relaisColisEventList     `json:"ListEvenements"`
}

type relaisColisEventList struct {
	Evenement []relaisColisEvent `json:"Evenement"`
}

type relaisColisEvent struct {
	Libelle string `json:"Libelle"`
	Date    string `json:"Date"`
	CodeJUS string `json:"CodeJUS"`
}

func parseRelaisColisResponse(data []byte) ([]model.TrackingEvent, error) {
	var resp relaisColisResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("relaiscolis: parse json: %w", err)
	}

	if resp.Colis == nil || resp.Colis.Colis == nil {
		return nil, fmt.Errorf("relaiscolis: no parcel data in response")
	}

	var events []model.TrackingEvent
	for _, e := range resp.Colis.Colis.ListEvenements.Evenement {
		ts, err := parseRelaisColisDate(e.Date)
		if err != nil {
			continue
		}

		events = append(events, model.TrackingEvent{
			Status:    mapRelaisColisStatus(e.CodeJUS, e.Libelle),
			Message:   e.Libelle,
			Timestamp: ts.UTC(),
		})
	}

	return events, nil
}

func parseRelaisColisDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"02/01/2006 15:04:05",
		"02/01/2006 15:04",
		"02/01/2006",
	}
	s = strings.TrimSpace(s)
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("relaiscolis: unknown date format: %q", s)
}

// mapRelaisColisStatus maps a Relais Colis event code or label to an internal ParcelStatus.
//
// Known codes:
//
//	LIV = Delivered (Livré)
func mapRelaisColisStatus(code, label string) model.ParcelStatus {
	upper := strings.ToUpper(strings.TrimSpace(code))

	if upper == "LIV" {
		return model.StatusDelivered
	}

	// Fall back to label-based mapping.
	lower := strings.ToLower(label)
	switch {
	// Out for delivery must be checked before delivered (both contain "livr").
	case strings.Contains(lower, "en cours de livraison") ||
		strings.Contains(lower, "mis en livraison"):
		return model.StatusOutForDelivery
	case strings.Contains(lower, "livr") ||
		strings.Contains(lower, "remis") ||
		strings.Contains(lower, "distribu"):
		return model.StatusDelivered
	case strings.Contains(lower, "pris en charge") ||
		strings.Contains(lower, "déposé") ||
		strings.Contains(lower, "enregistr"):
		return model.StatusInfoReceived
	case strings.Contains(lower, "retour") ||
		strings.Contains(lower, "refus") ||
		strings.Contains(lower, "anomalie"):
		return model.StatusFailed
	default:
		return model.StatusInTransit
	}
}
