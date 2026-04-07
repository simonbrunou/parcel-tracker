package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// MockTracker generates realistic fake tracking events for demos.
type MockTracker struct{}

func (t *MockTracker) Code() model.CarrierCode { return model.CarrierMock }
func (t *MockTracker) Name() string             { return "Mock (Demo)" }

func (t *MockTracker) Track(_ context.Context, trackingNumber string) ([]model.TrackingEvent, error) {
	now := time.Now().UTC()
	events := []model.TrackingEvent{
		{
			Status:    model.StatusInfoReceived,
			Message:   fmt.Sprintf("Shipment information received for %s", trackingNumber),
			Location:  "Origin facility",
			Timestamp: now.Add(-72 * time.Hour),
		},
		{
			Status:    model.StatusInTransit,
			Message:   "Package picked up by carrier",
			Location:  "Distribution center",
			Timestamp: now.Add(-48 * time.Hour),
		},
		{
			Status:    model.StatusInTransit,
			Message:   "Package in transit to destination",
			Location:  "Regional hub",
			Timestamp: now.Add(-24 * time.Hour),
		},
		{
			Status:    model.StatusOutForDelivery,
			Message:   "Out for delivery",
			Location:  "Local delivery facility",
			Timestamp: now.Add(-4 * time.Hour),
		},
		{
			Status:    model.StatusDelivered,
			Message:   "Delivered - Left at front door",
			Location:  "Destination",
			Timestamp: now.Add(-1 * time.Hour),
		},
	}
	return events, nil
}
