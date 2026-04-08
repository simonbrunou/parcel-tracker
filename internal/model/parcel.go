package model

import "time"

type ParcelStatus string

const (
	StatusUnknown      ParcelStatus = "unknown"
	StatusInfoReceived ParcelStatus = "info_received"
	StatusInTransit    ParcelStatus = "in_transit"
	StatusOutForDelivery ParcelStatus = "out_for_delivery"
	StatusDelivered    ParcelStatus = "delivered"
	StatusFailed       ParcelStatus = "failed"
	StatusExpired      ParcelStatus = "expired"
)

type CarrierCode string

const (
	CarrierManual  CarrierCode = "manual"
	CarrierMock    CarrierCode = "mock"
	CarrierUSPS    CarrierCode = "usps"
	CarrierFedEx   CarrierCode = "fedex"
	CarrierUPS     CarrierCode = "ups"
	CarrierDHL     CarrierCode = "dhl"
	CarrierPostNL  CarrierCode = "postnl"
	CarrierColissimo   CarrierCode = "colissimo"
	CarrierChronopost  CarrierCode = "chronopost"
	CarrierLaPoste     CarrierCode = "laposte"
)

type Parcel struct {
	ID             string       `json:"id"`
	TrackingNumber string       `json:"tracking_number"`
	Carrier        CarrierCode  `json:"carrier"`
	Name           string       `json:"name"`
	Notes          string       `json:"notes,omitempty"`
	Status         ParcelStatus `json:"status"`
	Archived       bool         `json:"archived"`
	LastCheck      *time.Time   `json:"last_check,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type TrackingEvent struct {
	ID        string       `json:"id"`
	ParcelID  string       `json:"parcel_id"`
	Status    ParcelStatus `json:"status"`
	Message   string       `json:"message"`
	Location  string       `json:"location,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
	CreatedAt time.Time    `json:"created_at"`
}
