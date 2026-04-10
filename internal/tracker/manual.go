package tracker

import (
	"context"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// ManualTracker is a no-op tracker for manually managed parcels.
type ManualTracker struct{}

func (t *ManualTracker) Code() model.CarrierCode { return model.CarrierManual }
func (t *ManualTracker) Name() string             { return "Manual" }

func (t *ManualTracker) Track(_ context.Context, _ string) (TrackResult, error) {
	return TrackResult{}, nil
}
