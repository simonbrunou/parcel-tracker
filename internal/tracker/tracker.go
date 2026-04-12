package tracker

import (
	"context"
	"os"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// TrackResult holds the outcome of a carrier tracking call.
// It includes the tracking events and an optional estimated delivery date.
type TrackResult struct {
	Events            []model.TrackingEvent
	EstimatedDelivery *time.Time
}

// Tracker defines how a carrier tracking provider works.
type Tracker interface {
	Code() model.CarrierCode
	Name() string
	Track(ctx context.Context, trackingNumber string) (TrackResult, error)
}

// CarrierInfo is returned to the frontend for the carrier dropdown.
type CarrierInfo struct {
	Code model.CarrierCode `json:"code"`
	Name string            `json:"name"`
}

// Registry holds all available trackers.
type Registry struct {
	trackers map[model.CarrierCode]Tracker
}

func NewRegistry() *Registry {
	r := &Registry{
		trackers: make(map[model.CarrierCode]Tracker),
	}

	laPosteAPIKey := os.Getenv("LAPOSTE_API_KEY")

	r.Register(&ManualTracker{})
	r.Register(&MockTracker{})
	r.Register(&ChronopostTracker{})
	if laPosteAPIKey != "" {
		r.Register(&LaPosteTracker{APIKey: laPosteAPIKey})
		r.Register(&ColissimoTracker{APIKey: laPosteAPIKey})
	}
	r.Register(&GLSTracker{})
	r.Register(&MondialRelayTracker{})
	r.Register(&DPDTracker{})
	r.Register(&ColisPriveTracker{})
	r.Register(&RelaisColisTracker{})
	r.Register(&VintedGoTracker{})
	return r
}

func (r *Registry) Register(t Tracker) {
	r.trackers[t.Code()] = t
}

func (r *Registry) Get(code model.CarrierCode) (Tracker, bool) {
	t, ok := r.trackers[code]
	return t, ok
}

func (r *Registry) Available() []CarrierInfo {
	var carriers []CarrierInfo
	for _, t := range r.trackers {
		carriers = append(carriers, CarrierInfo{
			Code: t.Code(),
			Name: t.Name(),
		})
	}
	return carriers
}
