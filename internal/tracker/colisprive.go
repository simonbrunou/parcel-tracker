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

const colisPriveTrackingURL = "https://www.colisprive.com/moncolis/pages/detailColis.aspx"

// ColisPriveTracker tracks parcels via the Colis Privé public tracking page.
type ColisPriveTracker struct {
	Client *http.Client
}

func (t *ColisPriveTracker) Code() model.CarrierCode { return model.CarrierColisPrive }
func (t *ColisPriveTracker) Name() string             { return "Colis Privé" }

func (t *ColisPriveTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *ColisPriveTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	url := colisPriveTrackingURL + "?numColis=" + trackingNumber

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResult{}, fmt.Errorf("colisprive: build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ParcelTracker/1.0)")
	req.Header.Set("Accept", "text/html")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("colisprive: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("colisprive: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return TrackResult{}, fmt.Errorf("colisprive: read response: %w", err)
	}

	return parseColisPriveHTML(body)
}

func parseColisPriveHTML(data []byte) (TrackResult, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return TrackResult{}, fmt.Errorf("colisprive: parse html: %w", err)
	}

	var result TrackResult

	// Look for estimated delivery date in the page text.
	var findDeliveryDate func(*html.Node)
	findDeliveryDate = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			lower := strings.ToLower(text)
			if strings.Contains(lower, "livraison prévue") || strings.Contains(lower, "livraison estimée") {
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

	// Colis Privé uses table rows with class "bandeauText" for tracking events.
	// Each row has two cells: date and label.
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" && hasClass(n, "bandeauText") {
			if ev, ok := parseColisPriveRow(n); ok {
				result.Events = append(result.Events, ev)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return result, nil
}

func parseColisPriveRow(n *html.Node) (model.TrackingEvent, bool) {
	var cells []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			text := strings.Join(collectTexts(c), " ")
			text = strings.TrimSpace(text)
			if text != "" {
				cells = append(cells, text)
			}
		}
	}

	if len(cells) < 2 {
		return model.TrackingEvent{}, false
	}

	dateStr := cells[0]
	label := cells[1]

	ts, err := parseColisPriveDate(dateStr)
	if err != nil {
		return model.TrackingEvent{}, false
	}

	return model.TrackingEvent{
		Status:    mapColisPriveStatus(label),
		Message:   strings.TrimSpace(label),
		Timestamp: ts.UTC(),
	}, true
}

func parseColisPriveDate(s string) (time.Time, error) {
	formats := []string{
		"02/01/2006 15:04",
		"02/01/2006 15h04",
		"02/01/2006",
		"2006-01-02 15:04:05",
		"2 January 2006 15:04",
	}
	s = strings.TrimSpace(s)
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("colisprive: unknown date format: %q", s)
}

// mapColisPriveStatus maps a Colis Privé event label to an internal ParcelStatus.
func mapColisPriveStatus(label string) model.ParcelStatus {
	lower := strings.ToLower(label)

	switch {
	// Out for delivery must be checked before delivered (both contain "livr").
	case strings.Contains(lower, "en cours de livraison") ||
		strings.Contains(lower, "mis en livraison") ||
		strings.Contains(lower, "en livraison") ||
		strings.Contains(lower, "chargé"):
		return model.StatusOutForDelivery
	// Failed must be checked before delivered ("non distribué" contains "distribu").
	case strings.Contains(lower, "non distribu") ||
		strings.Contains(lower, "retour") ||
		strings.Contains(lower, "refus") ||
		strings.Contains(lower, "anomalie") ||
		strings.Contains(lower, "avis de passage"):
		return model.StatusFailed
	case strings.Contains(lower, "livr") ||
		strings.Contains(lower, "remis") ||
		strings.Contains(lower, "distribu"):
		return model.StatusDelivered
	case strings.Contains(lower, "pris en charge") ||
		strings.Contains(lower, "réceptionné") ||
		strings.Contains(lower, "enregistr") ||
		strings.Contains(lower, "dédouanement"):
		return model.StatusInfoReceived
	default:
		return model.StatusInTransit
	}
}
