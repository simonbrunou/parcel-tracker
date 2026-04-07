package tracker

import (
	"context"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

// Tracker defines how a carrier tracking provider works.
type Tracker interface {
	Code() model.CarrierCode
	Name() string
	Track(ctx context.Context, trackingNumber string) ([]model.TrackingEvent, error)
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
	r.Register(&ManualTracker{})
	r.Register(&MockTracker{})
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
