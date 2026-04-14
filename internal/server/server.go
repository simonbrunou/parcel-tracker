package server

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
	"github.com/simonbrunou/parcel-tracker/internal/handler"
)

func New(ctx context.Context, h *handler.Handler, a *auth.Auth, distFS fs.FS, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(SecurityHeaders)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/ping"))

	rl := NewRateLimiter(10, time.Minute)
	rl.StartCleanup(ctx, 5*time.Minute)

	r.Route("/api", func(r chi.Router) {
		// Rate-limited auth endpoints
		r.Group(func(r chi.Router) {
			r.Use(rl.Middleware)
			r.Post("/auth/setup", h.Setup)
			r.Post("/auth/login", h.Login)
		})

		r.Post("/auth/logout", h.Logout)
		r.Get("/auth/check", h.CheckAuth)

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			handler.HealthCheck(w, r, h)
		})

		// VAPID key is public (needed before push subscription)
		r.Get("/notifications/vapid-key", h.GetVAPIDKey)

		r.Group(func(r chi.Router) {
			r.Use(a.Middleware)

			r.Route("/parcels", func(r chi.Router) {
				r.Get("/", h.ListParcels)
				r.Post("/", h.CreateParcel)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.GetParcel)
					r.Put("/", h.UpdateParcel)
					r.Delete("/", h.DeleteParcel)
					r.Post("/refresh", h.RefreshParcel)

					r.Get("/events", h.ListEvents)
					r.Post("/events", h.CreateEvent)
					r.Delete("/events/{eventID}", h.DeleteEvent)
				})
			})

			r.Post("/notifications/subscribe", h.Subscribe)
			r.Delete("/notifications/subscribe", h.Unsubscribe)
		})
	})

	r.Handle("/*", SPAHandler(distFS))

	return r
}
