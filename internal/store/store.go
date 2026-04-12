package store

import (
	"context"
	"time"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

type PushSubscription struct {
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
}

type ParcelFilter struct {
	Status   model.ParcelStatus
	Archived *bool
	Search   string
	Page     int
	PageSize int
}

type PaginatedParcels struct {
	Data     []model.Parcel `json:"data"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type Store interface {
	// Parcels
	ListParcels(ctx context.Context, filter ParcelFilter) (PaginatedParcels, error)
	ListActiveParcels(ctx context.Context) ([]model.Parcel, error)
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

	// Push Subscriptions
	ListPushSubscriptions(ctx context.Context) ([]PushSubscription, error)
	CreatePushSubscription(ctx context.Context, sub PushSubscription) (PushSubscription, error)
	DeletePushSubscription(ctx context.Context, endpoint string) error

	Close() error
}
