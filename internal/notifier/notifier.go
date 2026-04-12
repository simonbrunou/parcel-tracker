package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
)

// Payload is the JSON sent in each push message.
type Payload struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ParcelID string `json:"parcel_id"`
	Status   string `json:"status"`
}

// Notifier sends Web Push notifications when parcel statuses change.
type Notifier struct {
	Store  store.Store
	Logger *slog.Logger
}

// EnsureVAPIDKeys generates VAPID keys if they don't exist yet.
func (n *Notifier) EnsureVAPIDKeys(ctx context.Context) error {
	pub, _ := n.Store.GetSetting(ctx, "vapid_public_key")
	if pub != "" {
		return nil
	}

	priv, pub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return fmt.Errorf("generate VAPID keys: %w", err)
	}
	if err := n.Store.SetSetting(ctx, "vapid_private_key", priv); err != nil {
		return err
	}
	return n.Store.SetSetting(ctx, "vapid_public_key", pub)
}

// VAPIDPublicKey returns the stored VAPID public key.
func (n *Notifier) VAPIDPublicKey(ctx context.Context) (string, error) {
	return n.Store.GetSetting(ctx, "vapid_public_key")
}

// NotifyNewEvents sends a push notification for a parcel that has new events.
func (n *Notifier) NotifyNewEvents(ctx context.Context, parcel model.Parcel, newEvents int) {
	subs, err := n.Store.ListPushSubscriptions(ctx)
	if err != nil {
		n.Logger.Error("notifier: failed to list subscriptions", "error", err)
		return
	}
	if len(subs) == 0 {
		return
	}

	privKey, _ := n.Store.GetSetting(ctx, "vapid_private_key")
	pubKey, _ := n.Store.GetSetting(ctx, "vapid_public_key")
	if privKey == "" || pubKey == "" {
		return
	}

	title := parcel.Name
	if title == "" {
		title = parcel.TrackingNumber
	}

	payload := Payload{
		Title:    title,
		Body:     statusMessage(parcel.Status, newEvents),
		ParcelID: parcel.ID,
		Status:   string(parcel.Status),
	}

	msg, _ := json.Marshal(payload)

	for _, sub := range subs {
		s := &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.P256dh,
				Auth:   sub.Auth,
			},
		}
		resp, err := webpush.SendNotification(msg, s, &webpush.Options{
			VAPIDPublicKey:  pubKey,
			VAPIDPrivateKey: privKey,
			Subscriber:      "mailto:noreply@parcel-tracker.local",
			TTL:             3600,
		})
		if err != nil {
			n.Logger.Warn("notifier: push failed", "endpoint", sub.Endpoint, "error", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusGone {
			n.Logger.Info("notifier: removing stale subscription", "endpoint", sub.Endpoint)
			n.Store.DeletePushSubscription(ctx, sub.Endpoint)
		}
	}
}

func statusMessage(status model.ParcelStatus, newEvents int) string {
	switch status {
	case model.StatusDelivered:
		return "Your parcel has been delivered!"
	case model.StatusOutForDelivery:
		return "Your parcel is out for delivery"
	case model.StatusInTransit:
		return "Your parcel is in transit"
	case model.StatusFailed:
		return "Delivery attempt failed"
	case model.StatusInfoReceived:
		return "Shipping info received"
	default:
		if newEvents == 1 {
			return "1 new tracking update"
		}
		return fmt.Sprintf("%d new tracking updates", newEvents)
	}
}
