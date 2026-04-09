package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/simonbrunou/parcel-tracker/internal/auth"
	"github.com/simonbrunou/parcel-tracker/internal/model"
	"github.com/simonbrunou/parcel-tracker/internal/store"
	"github.com/simonbrunou/parcel-tracker/internal/tracker"
)

type testEnv struct {
	handler *Handler
	store   store.Store
	auth    *auth.Auth
	token   string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	s, err := store.NewSQLiteStore(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	a := auth.New(s)
	ctx := context.Background()
	a.Setup(ctx, "testpass")
	token, _ := a.Login(ctx, "testpass")

	reg := tracker.NewRegistry()

	h := &Handler{
		Store:   s,
		Auth:    a,
		Tracker: reg,
		Logger:  slog.Default(),
	}

	return &testEnv{handler: h, store: s, auth: a, token: token}
}

func (te *testEnv) authRequest(method, path string, body any) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	r := httptest.NewRequest(method, path, &buf)
	r.Header.Set("Content-Type", "application/json")
	r.AddCookie(&http.Cookie{Name: "session", Value: te.token})
	return r, httptest.NewRecorder()
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func withChiParams(r *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// --- Auth Handler Tests ---

func TestLoginSuccess(t *testing.T) {
	env := newTestEnv(t)

	body := bytes.NewBufferString(`{"password":"testpass"}`)
	r := httptest.NewRequest("POST", "/api/auth/login", body)
	w := httptest.NewRecorder()

	env.handler.Login(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected session cookie to be set")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	env := newTestEnv(t)

	body := bytes.NewBufferString(`{"password":"wrong"}`)
	r := httptest.NewRequest("POST", "/api/auth/login", body)
	w := httptest.NewRecorder()

	env.handler.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLoginBadJSON(t *testing.T) {
	env := newTestEnv(t)

	body := bytes.NewBufferString(`not json`)
	r := httptest.NewRequest("POST", "/api/auth/login", body)
	w := httptest.NewRecorder()

	env.handler.Login(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogout(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	env.handler.Logout(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	for _, c := range w.Result().Cookies() {
		if c.Name == "session" && c.MaxAge != -1 {
			t.Error("expected session cookie to be cleared")
		}
	}
}

func TestCheckAuthAuthenticated(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("GET", "/api/auth/check", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: env.token})
	w := httptest.NewRecorder()

	env.handler.CheckAuth(w, r)

	var resp map[string]bool
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp["authenticated"] {
		t.Error("expected authenticated=true")
	}
	if !resp["configured"] {
		t.Error("expected configured=true")
	}
}

func TestCheckAuthNotAuthenticated(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("GET", "/api/auth/check", nil)
	w := httptest.NewRecorder()

	env.handler.CheckAuth(w, r)

	var resp map[string]bool
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["authenticated"] {
		t.Error("expected authenticated=false")
	}
	if !resp["configured"] {
		t.Error("expected configured=true (password was set)")
	}
}

func TestSetupAlreadyConfigured(t *testing.T) {
	env := newTestEnv(t)

	body := bytes.NewBufferString(`{"password":"newpass"}`)
	r := httptest.NewRequest("POST", "/api/auth/setup", body)
	w := httptest.NewRecorder()

	env.handler.Setup(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestSetupPasswordTooShort(t *testing.T) {
	// Create a fresh env without setup
	s, _ := store.NewSQLiteStore(t.TempDir() + "/test.db")
	t.Cleanup(func() { s.Close() })
	a := auth.New(s)
	h := &Handler{Store: s, Auth: a, Tracker: tracker.NewRegistry(), Logger: slog.Default()}

	body := bytes.NewBufferString(`{"password":"ab"}`)
	r := httptest.NewRequest("POST", "/api/auth/setup", body)
	w := httptest.NewRecorder()

	h.Setup(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSetupSuccess(t *testing.T) {
	s, _ := store.NewSQLiteStore(t.TempDir() + "/test.db")
	t.Cleanup(func() { s.Close() })
	a := auth.New(s)
	h := &Handler{Store: s, Auth: a, Tracker: tracker.NewRegistry(), Logger: slog.Default()}

	body := bytes.NewBufferString(`{"password":"goodpass"}`)
	r := httptest.NewRequest("POST", "/api/auth/setup", body)
	w := httptest.NewRecorder()

	h.Setup(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Should have set session cookie
	found := false
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected session cookie after setup")
	}
}

// --- Parcel Handler Tests ---

func TestCreateParcel(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("POST", "/api/parcels", map[string]string{
		"tracking_number": "TRACK123",
		"carrier":         "manual",
		"name":            "My Package",
	})
	env.handler.CreateParcel(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var parcel model.Parcel
	json.NewDecoder(w.Body).Decode(&parcel)
	if parcel.TrackingNumber != "TRACK123" {
		t.Errorf("expected tracking_number TRACK123, got %s", parcel.TrackingNumber)
	}
	if parcel.Carrier != model.CarrierManual {
		t.Errorf("expected carrier manual, got %s", parcel.Carrier)
	}
}

func TestCreateParcelMissingTrackingNumber(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("POST", "/api/parcels", map[string]string{
		"carrier": "manual",
		"name":    "My Package",
	})
	env.handler.CreateParcel(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateParcelDefaultsToManual(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("POST", "/api/parcels", map[string]string{
		"tracking_number": "TRACK123",
	})
	env.handler.CreateParcel(w, r)

	var parcel model.Parcel
	json.NewDecoder(w.Body).Decode(&parcel)
	if parcel.Carrier != model.CarrierManual {
		t.Errorf("expected default carrier manual, got %s", parcel.Carrier)
	}
}

func TestListParcels(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "A", Carrier: model.CarrierManual})
	env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "B", Carrier: model.CarrierManual})

	r, w := env.authRequest("GET", "/api/parcels", nil)
	env.handler.ListParcels(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var parcels []model.Parcel
	json.NewDecoder(w.Body).Decode(&parcels)
	if len(parcels) != 2 {
		t.Errorf("expected 2 parcels, got %d", len(parcels))
	}
}

func TestListParcelsWithStatusFilter(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "A", Carrier: model.CarrierManual, Status: model.StatusInTransit})
	env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "B", Carrier: model.CarrierManual, Status: model.StatusDelivered})

	r, w := env.authRequest("GET", "/api/parcels?status=in_transit", nil)
	env.handler.ListParcels(w, r)

	var parcels []model.Parcel
	json.NewDecoder(w.Body).Decode(&parcels)
	if len(parcels) != 1 {
		t.Errorf("expected 1 in-transit parcel, got %d", len(parcels))
	}
}

func TestGetParcel(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual, Name: "Test"})

	r, w := env.authRequest("GET", "/api/parcels/"+p.ID, nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.GetParcel(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var parcel model.Parcel
	json.NewDecoder(w.Body).Decode(&parcel)
	if parcel.Name != "Test" {
		t.Errorf("expected name Test, got %s", parcel.Name)
	}
}

func TestGetParcelNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("GET", "/api/parcels/nonexistent", nil)
	r = withChiParam(r, "id", "nonexistent")
	env.handler.GetParcel(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateParcel(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC",
		Carrier:        model.CarrierManual,
		Name:           "Original",
	})

	r, w := env.authRequest("PUT", "/api/parcels/"+p.ID, map[string]any{
		"tracking_number": "ABC",
		"carrier":         "manual",
		"name":            "Updated",
		"status":          "in_transit",
	})
	r = withChiParam(r, "id", p.ID)
	env.handler.UpdateParcel(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var parcel model.Parcel
	json.NewDecoder(w.Body).Decode(&parcel)
	if parcel.Name != "Updated" {
		t.Errorf("expected name Updated, got %s", parcel.Name)
	}
}

func TestUpdateParcelNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("PUT", "/api/parcels/nonexistent", map[string]any{
		"tracking_number": "ABC",
		"carrier":         "manual",
		"name":            "Test",
	})
	r = withChiParam(r, "id", "nonexistent")
	env.handler.UpdateParcel(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteParcel(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	r, w := env.authRequest("DELETE", "/api/parcels/"+p.ID, nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.DeleteParcel(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestDeleteParcelNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("DELETE", "/api/parcels/nonexistent", nil)
	r = withChiParam(r, "id", "nonexistent")
	env.handler.DeleteParcel(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestRefreshParcelManualCarrier(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC",
		Carrier:        model.CarrierManual,
	})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/refresh", nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.RefreshParcel(w, r)

	// Manual tracker returns nil events, so refresh should succeed
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRefreshParcelNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("POST", "/api/parcels/nonexistent/refresh", nil)
	r = withChiParam(r, "id", "nonexistent")
	env.handler.RefreshParcel(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestRefreshParcelUnsupportedCarrier(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "ABC",
		Carrier:        "nonexistent_carrier",
	})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/refresh", nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.RefreshParcel(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRefreshParcelWithMockTracker(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{
		TrackingNumber: "MOCK1",
		Carrier:        model.CarrierMock,
	})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/refresh", nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.RefreshParcel(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify events were created
	events, _ := env.store.ListEvents(ctx, p.ID)
	if len(events) != 5 {
		t.Errorf("expected 5 events from mock tracker, got %d", len(events))
	}

	// Verify dedup: second refresh should not create duplicates
	r2, w2 := env.authRequest("POST", "/api/parcels/"+p.ID+"/refresh", nil)
	r2 = withChiParam(r2, "id", p.ID)
	env.handler.RefreshParcel(w2, r2)

	events, _ = env.store.ListEvents(ctx, p.ID)
	if len(events) != 5 {
		t.Errorf("expected 5 events after dedup, got %d", len(events))
	}
}

// --- Event Handler Tests ---

func TestListEvents(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	env.store.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "event1"})
	env.store.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "event2"})

	r, w := env.authRequest("GET", "/api/parcels/"+p.ID+"/events", nil)
	r = withChiParam(r, "id", p.ID)
	env.handler.ListEvents(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var events []model.TrackingEvent
	json.NewDecoder(w.Body).Decode(&events)
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestCreateEvent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/events", map[string]string{
		"message":  "Package received",
		"status":   "in_transit",
		"location": "Warehouse",
	})
	r = withChiParam(r, "id", p.ID)
	env.handler.CreateEvent(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var event model.TrackingEvent
	json.NewDecoder(w.Body).Decode(&event)
	if event.Message != "Package received" {
		t.Errorf("expected message 'Package received', got %q", event.Message)
	}
}

func TestCreateEventMissingMessage(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/events", map[string]string{
		"status": "in_transit",
	})
	r = withChiParam(r, "id", p.ID)
	env.handler.CreateEvent(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateEventDefaultStatus(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})

	r, w := env.authRequest("POST", "/api/parcels/"+p.ID+"/events", map[string]string{
		"message": "Something happened",
	})
	r = withChiParam(r, "id", p.ID)
	env.handler.CreateEvent(w, r)

	var event model.TrackingEvent
	json.NewDecoder(w.Body).Decode(&event)
	if event.Status != model.StatusUnknown {
		t.Errorf("expected default status %q, got %q", model.StatusUnknown, event.Status)
	}
}

func TestCreateEventParcelNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("POST", "/api/parcels/nonexistent/events", map[string]string{
		"message": "test",
	})
	r = withChiParam(r, "id", "nonexistent")
	env.handler.CreateEvent(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteEvent(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	p, _ := env.store.CreateParcel(ctx, model.Parcel{TrackingNumber: "ABC", Carrier: model.CarrierManual})
	e, _ := env.store.CreateEvent(ctx, model.TrackingEvent{ParcelID: p.ID, Message: "test"})

	r, w := env.authRequest("DELETE", "/api/parcels/"+p.ID+"/events/"+e.ID, nil)
	r = withChiParams(r, map[string]string{"id": p.ID, "eventID": e.ID})
	env.handler.DeleteEvent(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	env := newTestEnv(t)

	r, w := env.authRequest("DELETE", "/api/parcels/x/events/nonexistent", nil)
	r = withChiParams(r, map[string]string{"id": "x", "eventID": "nonexistent"})
	env.handler.DeleteEvent(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- Health Handler Tests ---

func TestHealthCheck(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	HealthCheck(w, r, env.handler)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %v", resp["status"])
	}
	carriers, ok := resp["carriers"].([]any)
	if !ok {
		t.Fatal("expected carriers array")
	}
	if len(carriers) < 1 {
		t.Error("expected at least 1 carrier")
	}
}

// --- Utility Function Tests ---

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"hello": "world"})

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["hello"] != "world" {
		t.Errorf("expected 'world', got %q", resp["hello"])
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "bad request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "bad request" {
		t.Errorf("expected 'bad request', got %q", resp["error"])
	}
}
