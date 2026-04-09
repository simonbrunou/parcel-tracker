package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/simonbrunou/parcel-tracker/internal/store"
)

func newTestAuth(t *testing.T) *Auth {
	t.Helper()
	s, err := store.NewSQLiteStore(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return New(s)
}

func TestSetupAndIsConfigured(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	if a.IsConfigured(ctx) {
		t.Fatal("expected not configured initially")
	}

	if err := a.Setup(ctx, "testpass"); err != nil {
		t.Fatalf("Setup: %v", err)
	}

	if !a.IsConfigured(ctx) {
		t.Fatal("expected configured after setup")
	}
}

func TestLoginSuccess(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "mypassword")

	token, err := a.Login(ctx, "mypassword")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "mypassword")

	_, err := a.Login(ctx, "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if err.Error() != "invalid password" {
		t.Errorf("expected 'invalid password', got %q", err.Error())
	}
}

func TestLoginNotConfigured(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	_, err := a.Login(ctx, "anything")
	if err == nil {
		t.Fatal("expected error when not configured")
	}
	if err.Error() != "not configured" {
		t.Errorf("expected 'not configured', got %q", err.Error())
	}
}

func TestVerifyValidToken(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "pass")
	token, _ := a.Login(ctx, "pass")

	if err := a.Verify(ctx, token); err != nil {
		t.Fatalf("Verify valid token: %v", err)
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "pass")

	if err := a.Verify(ctx, "garbage-token"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestVerifyTamperedToken(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "pass")
	token, _ := a.Login(ctx, "pass")

	// Tamper with the token by flipping a character
	tampered := token[:len(token)-1] + "X"
	if err := a.Verify(ctx, tampered); err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestExtractTokenFromCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "my-jwt-token"})

	got := ExtractToken(r)
	if got != "my-jwt-token" {
		t.Errorf("expected token from cookie, got %q", got)
	}
}

func TestExtractTokenFromAuthorizationHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer my-jwt-token")

	got := ExtractToken(r)
	if got != "my-jwt-token" {
		t.Errorf("expected token from header, got %q", got)
	}
}

func TestExtractTokenCookieTakesPrecedence(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "cookie-token"})
	r.Header.Set("Authorization", "Bearer header-token")

	got := ExtractToken(r)
	if got != "cookie-token" {
		t.Errorf("expected cookie token to take precedence, got %q", got)
	}
}

func TestExtractTokenEmpty(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	got := ExtractToken(r)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSetSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	SetSessionCookie(w, "test-token")

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != "session" {
		t.Errorf("expected cookie name 'session', got %q", c.Name)
	}
	if c.Value != "test-token" {
		t.Errorf("expected cookie value 'test-token', got %q", c.Value)
	}
	if !c.HttpOnly {
		t.Error("expected HttpOnly cookie")
	}
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("expected SameSiteStrict, got %v", c.SameSite)
	}
	if c.MaxAge <= 0 {
		t.Error("expected positive MaxAge")
	}
}

func TestClearSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ClearSessionCookie(w)

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}
	if cookies[0].MaxAge != -1 {
		t.Errorf("expected MaxAge=-1, got %d", cookies[0].MaxAge)
	}
}

func TestMiddlewareAllowsValidToken(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	a.Setup(ctx, "pass")
	token, _ := a.Login(ctx, "pass")

	called := false
	handler := a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: token})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("expected handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMiddlewareRejectsNoToken(t *testing.T) {
	a := newTestAuth(t)

	handler := a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()
	a.Setup(ctx, "pass")

	handler := a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "invalid-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTSecretIsPersistent(t *testing.T) {
	a := newTestAuth(t)
	ctx := context.Background()

	// Setup generates and stores a JWT secret
	a.Setup(ctx, "pass")
	token1, _ := a.Login(ctx, "pass")

	// A second login should use the same secret
	token2, _ := a.Login(ctx, "pass")

	// Both tokens should be verifiable
	if err := a.Verify(ctx, token1); err != nil {
		t.Errorf("first token should still be valid: %v", err)
	}
	if err := a.Verify(ctx, token2); err != nil {
		t.Errorf("second token should be valid: %v", err)
	}
}
