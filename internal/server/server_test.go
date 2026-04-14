package server

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
	"github.com/simonbrunou/parcel-tracker/internal/handler"
	"github.com/simonbrunou/parcel-tracker/internal/store"
	"github.com/simonbrunou/parcel-tracker/internal/tracker"
)

func newTestServer(t *testing.T, distFS fs.FS) http.Handler {
	t.Helper()
	s, err := store.NewSQLiteStore(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	a := auth.New(s)
	h := &handler.Handler{
		Store:   s,
		Auth:    a,
		Tracker: tracker.NewRegistry(),
		Logger:  slog.Default(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return New(ctx, h, a, distFS, slog.Default())
}

func TestPingEndpoint(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
	srv := newTestServer(t, distFS)

	r := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
	srv := newTestServer(t, distFS)

	r := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestProtectedEndpointRequiresAuth(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
	srv := newTestServer(t, distFS)

	r := httptest.NewRequest("GET", "/api/parcels/", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestSPAFallback(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>SPA</html>")},
	}
	srv := newTestServer(t, distFS)

	r := httptest.NewRequest("GET", "/some/unknown/path", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (SPA fallback), got %d", w.Code)
	}
}

func TestSPAHandlerServesStaticFiles(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html":        &fstest.MapFile{Data: []byte("<html></html>")},
		"assets/app.js":     &fstest.MapFile{Data: []byte("console.log('app')")},
	}

	handler := SPAHandler(distFS)

	// Serves existing file
	r := httptest.NewRequest("GET", "/assets/app.js", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for static file, got %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
		t.Errorf("expected aggressive cache for assets, got %q", cc)
	}
}

func TestSPAHandlerFallbackToIndex(t *testing.T) {
	distFS := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>SPA</html>")},
	}

	handler := SPAHandler(distFS)

	r := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (SPA fallback), got %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("expected no-cache for fallback, got %q", cc)
	}
}
