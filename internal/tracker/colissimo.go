package tracker

import (
	"context"
	"net/http"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// ColissimoTracker tracks Colissimo parcels via the La Poste Suivi v2 API.
// Colissimo is La Poste's parcel delivery service; both use the same tracking API.
type ColissimoTracker struct {
	APIKey string
	Client *http.Client
}

func (t *ColissimoTracker) Code() model.CarrierCode { return model.CarrierColissimo }
func (t *ColissimoTracker) Name() string             { return "Colissimo" }

func (t *ColissimoTracker) Track(ctx context.Context, trackingNumber string) ([]model.TrackingEvent, error) {
	delegate := &LaPosteTracker{
		APIKey: t.APIKey,
		Client: t.Client,
	}
	return delegate.Track(ctx, trackingNumber)
}
