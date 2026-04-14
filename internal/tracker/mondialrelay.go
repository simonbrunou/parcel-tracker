package tracker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

const mondialRelayTrackingURL = "https://www.mondialrelay.fr/suivi-de-colis"

// MondialRelayTracker tracks parcels via the Mondial Relay public tracking page.
//
// Tracking numbers must be in the format "{expeditionNumber}-{postalCode}"
// (e.g. "12345678-75001"). The expedition number is 8, 10, or 12 digits
// and the postal code is 5 digits.
type MondialRelayTracker struct {
	Client *http.Client
}

func (t *MondialRelayTracker) Code() model.CarrierCode { return model.CarrierMondialRelay }
func (t *MondialRelayTracker) Name() string             { return "Mondial Relay" }

func (t *MondialRelayTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *MondialRelayTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	parts := strings.SplitN(trackingNumber, "-", 2)
	if len(parts) != 2 {
		return TrackResult{}, fmt.Errorf("mondialrelay: tracking number must be in format 'expeditionNumber-postalCode' (e.g. '12345678-75001')")
	}
	expeditionNumber := parts[0]
	postalCode := parts[1]

	url := fmt.Sprintf("%s?numeroExpedition=%s&codePostal=%s", mondialRelayTrackingURL, expeditionNumber, postalCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResult{}, fmt.Errorf("mondialrelay: build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ParcelTracker/1.0)")
	req.Header.Set("Accept", "text/html")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("mondialrelay: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("mondialrelay: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return TrackResult{}, fmt.Errorf("mondialrelay: read response: %w", err)
	}

	return parseMondialRelayHTML(body)
}

func parseMondialRelayHTML(data []byte) (TrackResult, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return TrackResult{}, fmt.Errorf("mondialrelay: parse html: %w", err)
	}

	var result TrackResult

	// Look for estimated delivery date in the page text.
	var findDeliveryDate func(*html.Node)
	findDeliveryDate = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			lower := strings.ToLower(text)
			if strings.Contains(lower, "livraison estimée") || strings.Contains(lower, "livraison prévue") {
				if t, err := extractDateFromText(text); err == nil {
					utc := t.UTC()
					result.EstimatedDelivery = &utc
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findDeliveryDate(c)
		}
	}
	findDeliveryDate(doc)

	// Mondial Relay uses div elements with class "infos-account" for the timeline.
	// Each block contains a date and sub-events with time and label.
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "infos-account") {
			blockEvents := parseMondialRelayBlock(n)
			result.Events = append(result.Events, blockEvents...)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	if len(result.Events) == 0 {
		// Try alternative structure: look for timeline items by other class patterns.
		var walkAlt func(*html.Node)
		walkAlt = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "timeline-item") {
				if ev, ok := parseMondialRelayTimelineItem(n); ok {
					result.Events = append(result.Events, ev)
				}
				return
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walkAlt(c)
			}
		}
		walkAlt(doc)
	}

	return result, nil
}

// extractDateFromText tries to find and parse a date from a text string
// containing patterns like "Livraison estimée le 03/06/2025".
func extractDateFromText(text string) (time.Time, error) {
	formats := []string{
		"02/01/2006",
		"2006-01-02",
	}
	// Try to find a date-like substring in the text.
	words := strings.Fields(text)
	for _, word := range words {
		clean := strings.Trim(word, ".,;:")
		for _, f := range formats {
			if t, err := time.Parse(f, clean); err == nil {
				return t, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("no date found in text: %q", text)
}

func parseMondialRelayBlock(n *html.Node) []model.TrackingEvent {
	var events []model.TrackingEvent
	texts := collectTexts(n)

	// The block has a date followed by pairs of (time, label).
	if len(texts) < 2 {
		return nil
	}

	dateStr := strings.TrimSpace(texts[0])
	for i := 1; i+1 < len(texts); i += 2 {
		timeStr := strings.TrimSpace(texts[i])
		label := strings.TrimSpace(texts[i+1])

		ts, err := parseMondialRelayDate(dateStr, timeStr)
		if err != nil {
			continue
		}

		events = append(events, model.TrackingEvent{
			Status:    mapMondialRelayStatus(label),
			Message:   label,
			Timestamp: ts.UTC(),
		})
	}

	return events
}

func parseMondialRelayTimelineItem(n *html.Node) (model.TrackingEvent, bool) {
	texts := collectTexts(n)
	if len(texts) < 2 {
		return model.TrackingEvent{}, false
	}

	dateStr := strings.TrimSpace(texts[0])
	label := strings.TrimSpace(texts[len(texts)-1])

	ts, err := time.Parse("02/01/2006 15:04", dateStr)
	if err != nil {
		ts, err = time.Parse("02/01/2006", dateStr)
		if err != nil {
			return model.TrackingEvent{}, false
		}
	}

	return model.TrackingEvent{
		Status:    mapMondialRelayStatus(label),
		Message:   label,
		Timestamp: ts.UTC(),
	}, true
}

func parseMondialRelayDate(date, timeStr string) (time.Time, error) {
	combined := date + " " + timeStr
	formats := []string{
		"02/01/2006 15:04",
		"02/01/2006 15h04",
		"2 January 2006 15:04",
		"02/01/2006",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, combined); err == nil {
			return t, nil
		}
	}
	// Try date alone.
	if t, err := time.Parse("02/01/2006", date); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("mondialrelay: unknown date format: %q %q", date, timeStr)
}

// mapMondialRelayStatus maps a Mondial Relay event label to an internal ParcelStatus.
func mapMondialRelayStatus(label string) model.ParcelStatus {
	lower := strings.ToLower(label)

	switch {
	// Out for delivery must be checked before delivered (both contain "livr").
	case strings.Contains(lower, "en cours de livraison") ||
		strings.Contains(lower, "mis en livraison"):
		return model.StatusOutForDelivery
	case strings.Contains(lower, "livr") ||
		strings.Contains(lower, "retir") ||
		strings.Contains(lower, "distribu"):
		return model.StatusDelivered
	case strings.Contains(lower, "pris en charge") ||
		strings.Contains(lower, "déposé") ||
		strings.Contains(lower, "collecté") ||
		strings.Contains(lower, "enregistr"):
		return model.StatusInfoReceived
	case strings.Contains(lower, "anomalie") ||
		strings.Contains(lower, "retour") ||
		strings.Contains(lower, "refus") ||
		strings.Contains(lower, "non distribuable"):
		return model.StatusFailed
	default:
		return model.StatusInTransit
	}
}

// HTML utility functions used by all scrapers.

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, c := range strings.Fields(attr.Val) {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func hasAttrContaining(n *html.Node, key, substr string) bool {
	for _, attr := range n.Attr {
		if attr.Key == key && strings.Contains(attr.Val, substr) {
			return true
		}
	}
	return false
}

func collectTexts(n *html.Node) []string {
	var texts []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			t := strings.TrimSpace(n.Data)
			if t != "" {
				texts = append(texts, t)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return texts
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
