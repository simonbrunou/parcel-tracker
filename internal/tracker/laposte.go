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

const laPosteTrackingURL = "https://api.laposte.fr/suivi/v2/idships/"

// LaPosteTracker tracks parcels via the La Poste Suivi v2 REST API.
// It covers La Poste letters, Colissimo, and Chronopost shipments.
// Requires an API key from https://developer.laposte.fr.
type LaPosteTracker struct {
	APIKey string
	Client *http.Client
}

func (t *LaPosteTracker) Code() model.CarrierCode { return model.CarrierLaPoste }
func (t *LaPosteTracker) Name() string             { return "La Poste" }

func (t *LaPosteTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *LaPosteTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	if t.APIKey == "" {
		return TrackResult{}, fmt.Errorf("laposte: LAPOSTE_API_KEY environment variable is not set")
	}

	url := laPosteTrackingURL + trackingNumber + "?lang=fr_FR"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResult{}, fmt.Errorf("laposte: build request: %w", err)
	}
	req.Header.Set("X-Okapi-Key", t.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("laposte: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackResult{}, fmt.Errorf("laposte: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("laposte: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return parseLaPosteResponse(body)
}

// La Poste JSON response structures.

type laPosteResponse struct {
	Shipment      *laPosteShipment `json:"shipment"`
	ReturnMessage string           `json:"returnMessage"`
}

type laPosteShipment struct {
	Product      string          `json:"product"`
	ContextData  *laPosteContext `json:"contextData"`
	DeliveryDate string          `json:"deliveryDate"`
	Event        []laPosteEvent  `json:"event"`
}

type laPosteContext struct {
	OriginCountry  string `json:"originCountry"`
	ArrivalCountry string `json:"arrivalCountry"`
}

type laPosteEvent struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Date  string `json:"date"`
}

func parseLaPosteResponse(data []byte) (TrackResult, error) {
	var resp laPosteResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return TrackResult{}, fmt.Errorf("laposte: parse json: %w", err)
	}

	if resp.Shipment == nil {
		msg := resp.ReturnMessage
		if msg == "" {
			msg = "no shipment data"
		}
		return TrackResult{}, fmt.Errorf("laposte: %s", msg)
	}

	var result TrackResult

	if resp.Shipment.DeliveryDate != "" {
		if t, err := parseLaPosteDate(resp.Shipment.DeliveryDate); err == nil {
			utc := t.UTC()
			result.EstimatedDelivery = &utc
		}
	}

	for _, e := range resp.Shipment.Event {
		ts, err := parseLaPosteDate(e.Date)
		if err != nil {
			continue
		}

		result.Events = append(result.Events, model.TrackingEvent{
			Status:    mapLaPosteStatus(e.Code),
			Message:   e.Label,
			Timestamp: ts.UTC(),
		})
	}

	return result, nil
}

func parseLaPosteDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05.000",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("laposte: unknown date format: %q", s)
}

// mapLaPosteStatus maps a La Poste event code to an internal ParcelStatus.
//
// La Poste event codes:
//
//	DR1       = Shipment information received
//	PC1, PC2  = Picked up / accepted
//	ET1-ET4   = In transit / being processed
//	EP1       = Awaiting presentation
//	DO1, DO2  = Customs entry/exit
//	DO3       = Held at customs
//	PB1       = Problem in progress
//	PB2       = Problem resolved
//	MD2       = Out for delivery
//	ND1       = Undeliverable
//	AG1       = Awaiting pickup at post office
//	RE1       = Returned to sender
//	DI1       = Delivered
//	DI2       = Delivered to sender (return)
func mapLaPosteStatus(code string) model.ParcelStatus {
	upper := strings.ToUpper(strings.TrimSpace(code))

	switch upper {
	// Delivered
	case "DI1", "DI2":
		return model.StatusDelivered
	// Out for delivery
	case "MD2":
		return model.StatusOutForDelivery
	// Info received / picked up
	case "DR1", "PC1", "PC2":
		return model.StatusInfoReceived
	// Failed / returned / undeliverable
	case "PB1", "ND1", "RE1":
		return model.StatusFailed
	// In transit (processing, customs, awaiting)
	default:
		return model.StatusInTransit
	}
}
