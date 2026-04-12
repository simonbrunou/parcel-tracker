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
	CarrierColissimo    CarrierCode = "colissimo"
	CarrierChronopost   CarrierCode = "chronopost"
	CarrierLaPoste      CarrierCode = "laposte"
	CarrierMondialRelay CarrierCode = "mondialrelay"
	CarrierGLS          CarrierCode = "gls"
	CarrierDPD          CarrierCode = "dpd"
	CarrierColisPrive   CarrierCode = "colisprive"
	CarrierRelaisColis  CarrierCode = "relaiscolis"
	CarrierVintedGo     CarrierCode = "vintedgo"
)

// IsValid reports whether the carrier code is a known carrier.
func (c CarrierCode) IsValid() bool {
	switch c {
	case CarrierManual, CarrierMock, CarrierUSPS, CarrierFedEx, CarrierUPS, CarrierDHL,
		CarrierPostNL, CarrierColissimo, CarrierChronopost, CarrierLaPoste,
		CarrierMondialRelay, CarrierGLS, CarrierDPD, CarrierColisPrive, CarrierRelaisColis,
		CarrierVintedGo:
		return true
	}
	return false
}

// IsValid reports whether the parcel status is a known status.
func (s ParcelStatus) IsValid() bool {
	switch s {
	case StatusUnknown, StatusInfoReceived, StatusInTransit, StatusOutForDelivery,
		StatusDelivered, StatusFailed, StatusExpired:
		return true
	}
	return false
}

type Parcel struct {
	ID             string       `json:"id"`
	TrackingNumber string       `json:"tracking_number"`
	Carrier        CarrierCode  `json:"carrier"`
	Name           string       `json:"name"`
	Notes          string       `json:"notes,omitempty"`
	Status         ParcelStatus `json:"status"`
	Archived          bool         `json:"archived"`
	EstimatedDelivery *time.Time   `json:"estimated_delivery,omitempty"`
	LastCheck         *time.Time   `json:"last_check,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
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
