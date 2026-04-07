package store

import (
	"context"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

type ParcelFilter struct {
	Status   model.ParcelStatus
	Archived *bool
	Search   string
}

type Store interface {
	// Parcels
	ListParcels(ctx context.Context, filter ParcelFilter) ([]model.Parcel, error)
	GetParcel(ctx context.Context, id string) (model.Parcel, error)
	CreateParcel(ctx context.Context, p model.Parcel) (model.Parcel, error)
	UpdateParcel(ctx context.Context, p model.Parcel) (model.Parcel, error)
	DeleteParcel(ctx context.Context, id string) error

	// Tracking Events
	ListEvents(ctx context.Context, parcelID string) ([]model.TrackingEvent, error)
	CreateEvent(ctx context.Context, e model.TrackingEvent) (model.TrackingEvent, error)
	DeleteEvent(ctx context.Context, id string) error

	// Settings
	GetSetting(ctx context.Context, key string) (string, error)
	SetSetting(ctx context.Context, key, value string) error

	Close() error
}
