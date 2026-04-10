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

const dpdTrackingURL = "https://trace.dpd.fr/fr/trace/"

// DPDTracker tracks parcels via the DPD France public tracking page.
type DPDTracker struct {
	Client *http.Client
}

func (t *DPDTracker) Code() model.CarrierCode { return model.CarrierDPD }
func (t *DPDTracker) Name() string             { return "DPD" }

func (t *DPDTracker) httpClient() *http.Client {
	if t.Client != nil {
		return t.Client
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (t *DPDTracker) Track(ctx context.Context, trackingNumber string) (TrackResult, error) {
	url := dpdTrackingURL + trackingNumber

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TrackResult{}, fmt.Errorf("dpd: build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ParcelTracker/1.0)")
	req.Header.Set("Accept", "text/html")

	resp, err := t.httpClient().Do(req)
	if err != nil {
		return TrackResult{}, fmt.Errorf("dpd: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TrackResult{}, fmt.Errorf("dpd: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackResult{}, fmt.Errorf("dpd: read response: %w", err)
	}

	return parseDPDHTML(body)
}

func parseDPDHTML(data []byte) (TrackResult, error) {
	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return TrackResult{}, fmt.Errorf("dpd: parse html: %w", err)
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

	// DPD uses table rows with id containing "ligneTableTrace" for tracking events.
	// Each row has cells: date, time, label, [location].
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" && hasAttrContaining(n, "id", "ligneTableTrace") {
			if ev, ok := parseDPDRow(n); ok {
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

func parseDPDRow(n *html.Node) (model.TrackingEvent, bool) {
	// Collect text from each <td> cell.
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

	// Expect at least: date, time, label.
	if len(cells) < 3 {
		return model.TrackingEvent{}, false
	}

	dateStr := cells[0]
	timeStr := cells[1]
	label := cells[2]
	location := ""
	if len(cells) >= 4 {
		location = cells[3]
	}

	ts, err := parseDPDDate(dateStr, timeStr)
	if err != nil {
		return model.TrackingEvent{}, false
	}

	return model.TrackingEvent{
		Status:    mapDPDStatus(label),
		Message:   cleanDPDLabel(label),
		Location:  location,
		Timestamp: ts.UTC(),
	}, true
}

func parseDPDDate(date, timeStr string) (time.Time, error) {
	combined := strings.TrimSpace(date) + " " + strings.TrimSpace(timeStr)
	formats := []string{
		"02/01/2006 15:04",
		"02/01/2006 15h04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, combined); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("dpd: unknown date format: %q", combined)
}

func cleanDPDLabel(label string) string {
	// Remove common DPD prefixes from labels.
	label = strings.ReplaceAll(label, "Predict vous informe : \n", "")
	label = strings.ReplaceAll(label, "Predict vous informe : ", "")
	label = strings.ReplaceAll(label, "Instruction :", "")
	return strings.TrimSpace(label)
}

// mapDPDStatus maps a DPD event label to an internal ParcelStatus.
func mapDPDStatus(label string) model.ParcelStatus {
	lower := strings.ToLower(label)

	switch {
	// Out for delivery must be checked before delivered (both contain "livr").
	case strings.Contains(lower, "en cours de livraison") ||
		strings.Contains(lower, "mis en livraison") ||
		strings.Contains(lower, "en livraison") ||
		strings.Contains(lower, "avec le chauffeur"):
		return model.StatusOutForDelivery
	case strings.Contains(lower, "livr") ||
		strings.Contains(lower, "remis") ||
		strings.Contains(lower, "delivered"):
		return model.StatusDelivered
	case strings.Contains(lower, "pris en charge") ||
		strings.Contains(lower, "enlev") ||
		strings.Contains(lower, "enregistr") ||
		strings.Contains(lower, "données de colis"):
		return model.StatusInfoReceived
	case strings.Contains(lower, "non distribu") ||
		strings.Contains(lower, "retour") ||
		strings.Contains(lower, "refus") ||
		strings.Contains(lower, "anomalie"):
		return model.StatusFailed
	default:
		return model.StatusInTransit
	}
}
