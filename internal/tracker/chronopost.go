package tracker

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const (
	chronopostTrackingURL = "https://ws.chronopost.fr/tracking-cxf/TrackingServiceWS"
	chronopostSOAPAction  = ""
)

// ChronopostTracker tracks parcels via the Chronopost SOAP tracking API.
type ChronopostTracker struct {
	Client *http.Client
}

func (t *ChronopostTracker) Code() model.CarrierCode { return model.CarrierChronopost }
func (t *ChronopostTracker) Name() string             { return "Chronopost" }

func (t *ChronopostTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *ChronopostTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	body := buildChronopostSOAPRequest(trackingNumber, "fr")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chronopostTrackingURL, bytes.NewReader(body))
	if err != nil {
		return TrackResult{}, fmt.Errorf("chronopost: build request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", chronopostSOAPAction)

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("chronopost: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackResult{}, fmt.Errorf("chronopost: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("chronopost: unexpected status %d", resp.StatusCode)
	}

	return parseChronopostResponse(respBody)
}

func buildChronopostSOAPRequest(trackingNumber, language string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	buf.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cxf="http://cxf.tracking.soap.chronopost.fr/">`)
	buf.WriteString(`<soapenv:Body>`)
	buf.WriteString(`<cxf:trackSkybillV2>`)
	buf.WriteString(`<language>`)
	buf.WriteString(language)
	buf.WriteString(`</language>`)
	buf.WriteString(`<skybillNumber>`)
	buf.WriteString(trackingNumber)
	buf.WriteString(`</skybillNumber>`)
	buf.WriteString(`</cxf:trackSkybillV2>`)
	buf.WriteString(`</soapenv:Body>`)
	buf.WriteString(`</soapenv:Envelope>`)
	return buf.Bytes()
}

// SOAP response XML structures.

type chronopostEnvelope struct {
	XMLName xml.Name        `xml:"Envelope"`
	Body    chronopostBody  `xml:"Body"`
}

type chronopostBody struct {
	Response chronopostTrackResponse `xml:"trackSkybillV2Response"`
	Fault    *chronopostFault        `xml:"Fault"`
}

type chronopostFault struct {
	FaultString string `xml:"faultstring"`
}

type chronopostTrackResponse struct {
	Return chronopostReturn `xml:"return"`
}

type chronopostReturn struct {
	ErrorCode              int                     `xml:"errorCode"`
	ErrorMessage           string                  `xml:"errorMessage"`
	EstimatedDeliveryDate  string                  `xml:"estimatedDeliveryDate"`
	ListEventInfoComp      chronopostEventInfoComp `xml:"listEventInfoComp"`
}

type chronopostEventInfoComp struct {
	SkybillNumber string              `xml:"skybillNumber"`
	Events        []chronopostEvent   `xml:"events"`
}

type chronopostEvent struct {
	Code        string `xml:"code"`
	Label       string `xml:"eventLabel"`
	Date        string `xml:"eventDate"`
	OfficeLabel string `xml:"officeLabel"`
	ZipCode     string `xml:"zipCode"`
}

func parseChronopostResponse(data []byte) (TrackResult, error) {
	var env chronopostEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return TrackResult{}, fmt.Errorf("chronopost: parse xml: %w", err)
	}

	if env.Body.Fault != nil {
		return TrackResult{}, fmt.Errorf("chronopost: SOAP fault: %s", env.Body.Fault.FaultString)
	}

	ret := env.Body.Response.Return
	if ret.ErrorCode != 0 {
		return TrackResult{}, fmt.Errorf("chronopost: error %d: %s", ret.ErrorCode, ret.ErrorMessage)
	}

	var result TrackResult

	if ret.EstimatedDeliveryDate != "" {
		if t, err := parseChronopostDate(ret.EstimatedDeliveryDate); err == nil {
			utc := t.UTC()
			result.EstimatedDelivery = &utc
		}
	}

	for _, e := range ret.ListEventInfoComp.Events {
		ts, err := parseChronopostDate(e.Date)
		if err != nil {
			continue
		}

		location := buildLocation(e.OfficeLabel, e.ZipCode)

		result.Events = append(result.Events, model.TrackingEvent{
			Status:    mapChronopostStatus(e.Code),
			Message:   e.Label,
			Location:  location,
			Timestamp: ts.UTC(),
		})
	}

	return result, nil
}

func parseChronopostDate(s string) (time.Time, error) {
	// Chronopost returns dates in various formats; try the most common ones.
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
	return time.Time{}, fmt.Errorf("chronopost: unknown date format: %q", s)
}

func buildLocation(site, zipCode string) string {
	site = strings.TrimSpace(site)
	zipCode = strings.TrimSpace(zipCode)
	switch {
	case site != "" && zipCode != "":
		return site + " (" + zipCode + ")"
	case site != "":
		return site
	case zipCode != "":
		return zipCode
	default:
		return ""
	}
}

// mapChronopostStatus maps a Chronopost event code to an internal ParcelStatus.
//
// Codes observed in the live tracking-cxf SOAP API (Chronopost returns them
// padded with trailing spaces, hence the TrimSpace):
//
//	DC = Preparation at sender (label created) — info received
//	EP = Picked up from sender                  — info received
//	PH = Package received by Chronopost         — info received
//	RG = Registered                              — info received
//	EC = Sorted at origin hub                    — in transit
//	TS = In transit between hubs                 — in transit
//	TA = Arrived at sorting hub                  — in transit
//	TI = In transit                              — in transit
//	SD = Sorted at destination depot             — in transit
//	IS = Informational (e.g. delivery scheduled) — in transit
//	LT = Out for delivery (livreur en tournée)   — out for delivery
//	CR = Courier en route                        — out for delivery
//	MD = Mise en distribution                    — out for delivery
//	LV = Delivered (livré)                       — delivered
//	RM = Remis au destinataire                   — delivered
//	D1,D2,D3... = Historical delivered variants  — delivered
//	LE = Return to sender                        — failed
//	RE = Refused / returned                      — failed
//	AR = Anomaly / issue                         — failed
//
// Unknown codes fall back to StatusInTransit.
func mapChronopostStatus(code string) model.ParcelStatus {
	upper := strings.ToUpper(strings.TrimSpace(code))

	switch {
	// Delivered. Match LV, RM, and legacy D + digit variants (D1, D2, ...).
	// We deliberately do NOT match every "D*" because DC means "preparation
	// at sender", not delivered.
	case strings.HasPrefix(upper, "LV"),
		strings.HasPrefix(upper, "RM"),
		len(upper) >= 2 && upper[0] == 'D' && upper[1] >= '0' && upper[1] <= '9':
		return model.StatusDelivered
	// Out for delivery.
	case strings.HasPrefix(upper, "LT"),
		strings.HasPrefix(upper, "CR"),
		strings.HasPrefix(upper, "MD"):
		return model.StatusOutForDelivery
	// Info received / picked up / label created at sender.
	case strings.HasPrefix(upper, "DC"),
		strings.HasPrefix(upper, "EP"),
		strings.HasPrefix(upper, "PH"),
		strings.HasPrefix(upper, "RG"):
		return model.StatusInfoReceived
	// Failed / returned / refused / anomaly.
	case strings.HasPrefix(upper, "LE"),
		strings.HasPrefix(upper, "RE"),
		strings.HasPrefix(upper, "AR"):
		return model.StatusFailed
	// In transit: EC, TS, TA, TI, SD, IS and any unknown code.
	default:
		return model.StatusInTransit
	}
}
