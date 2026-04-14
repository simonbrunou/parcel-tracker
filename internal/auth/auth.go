package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/simonbrunou/parcel-tracker/internal/store"
)

const (
	settingPasswordHash = "password_hash"
	settingJWTSecret    = "jwt_secret"
	cookieName          = "session"
	tokenLifetime       = 30 * 24 * time.Hour
)

type Auth struct {
	store store.Store
	mu    sync.Mutex
}

func New(s store.Store) *Auth {
	return &Auth{store: s}
}

// Setup hashes the password and stores it, generating a JWT secret if needed.
func (a *Auth) Setup(ctx context.Context, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := a.store.SetSetting(ctx, settingPasswordHash, string(hash)); err != nil {
		return err
	}
	// Ensure JWT secret exists
	_, err = a.jwtSecret(ctx)
	return err
}

// IsConfigured returns true if a password has been set.
func (a *Auth) IsConfigured(ctx context.Context) bool {
	hash, _ := a.store.GetSetting(ctx, settingPasswordHash)
	return hash != ""
}

// Login verifies the password and returns a JWT token.
func (a *Auth) Login(ctx context.Context, password string) (string, error) {
	hash, err := a.store.GetSetting(ctx, settingPasswordHash)
	if err != nil || hash == "" {
		return "", errors.New("not configured")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	secret, err := a.jwtSecret(ctx)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(tokenLifetime).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(secret)
}

// Verify checks if the token is valid.
func (a *Auth) Verify(ctx context.Context, tokenStr string) error {
	secret, err := a.jwtSecret(ctx)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// SetSessionCookie sets the JWT as an HttpOnly cookie.
// When secure is true, the cookie is only sent over HTTPS.
func SetSessionCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(tokenLifetime.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// ExtractToken gets the token from the cookie or Authorization header.
func ExtractToken(r *http.Request) string {
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

func (a *Auth) jwtSecret(ctx context.Context) ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	secretHex, err := a.store.GetSetting(ctx, settingJWTSecret)
	if err != nil {
		return nil, err
	}
	if secretHex != "" {
		return hex.DecodeString(secretHex)
	}

	// Generate new secret
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	secretHex = hex.EncodeToString(secret)
	if err := a.store.SetSetting(ctx, settingJWTSecret, secretHex); err != nil {
		return nil, err
	}
	return secret, nil
}
