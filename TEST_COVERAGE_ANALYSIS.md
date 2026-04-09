# Test Coverage Analysis

_Generated: 2026-04-09_

## Current State

Overall test coverage is low. Only `internal/tracker` has tests (53.1% — parsing helpers only). All other Go packages and the entire frontend have 0% coverage.

| Package | Coverage | Test Files | Notes |
|---|---|---|---|
| `internal/tracker` | **53.1%** | 8 files, 49 tests | Parsing, status mapping, dates |
| `internal/auth` | **0%** | 0 | Security-critical |
| `internal/handler` | **0%** | 0 | All HTTP handlers |
| `internal/store` | **0%** | 0 | SQLite data layer |
| `internal/server` | **0%** | 0 | Router, middleware, SPA |
| `internal/config` | **0%** | 0 | Env var parsing |
| `internal/model` | N/A | 0 | Constants only |
| `web/` (frontend) | **0%** | 0 | No test framework configured |

### What's Tested Well

The tracker package tests follow good patterns:
- `httptest.NewServer` for mocking carrier APIs
- Table-driven tests for status mapping and date parsing
- Error case coverage (missing API keys, malformed responses, SOAP faults)
- Tests for all 8 real carrier implementations' parsing logic

### What's Not Tested

Within the tracker package, `Track()` methods, `Registry`, `Worker`, and `EventKey()` have 0% coverage.

---

## Recommended Improvements (Priority Order)

### 1. `internal/auth` — High Priority (Security-Critical)

**Why:** Handles authentication for the entire application. Bugs here mean unauthorized access.

**Tests to add:**
- `Setup()` — hashes password, stores it, prevents duplicate setup (409)
- `Login()` — correct password returns JWT, wrong password fails
- `Verify()` — valid token passes, expired token rejected, tampered token rejected
- `AuthMiddleware()` — missing token returns 401, valid token passes request through
- `ExtractToken()` — reads cookie first, falls back to Authorization Bearer header
- `SetSessionCookie()` / `ClearSessionCookie()` — correct cookie attributes (HttpOnly, SameSite, MaxAge)

**Approach:** Use an in-memory SQLite store (or mock Store interface) to avoid filesystem dependencies. Consider a lower bcrypt cost for test speed.

**Bug found — `jwtSecret()` race condition** (`auth.go:131`): Multiple goroutines calling `jwtSecret()` simultaneously could both generate and save different secrets. Needs a mutex or `sync.Once`.

### 2. `internal/store` — High Priority (Data Integrity)

**Why:** All data flows through here. Bugs mean data loss or corruption.

**Tests to add:**
- `CreateParcel` / `GetParcel` / `UpdateParcel` / `DeleteParcel` — CRUD round-trip
- `ListParcels` — filter by status, search (case-insensitive LIKE), archived flag, ordering by updated_at
- `CreateEvent` — verify the side-effect that updates parcel status
- `DeleteEvent` — verify 404 behavior
- `GetSetting` / `SetSetting` — round-trip, missing key returns empty string
- `migrate()` — schema creation with tables, indexes, foreign keys

**Approach:** Use `:memory:` SQLite for fast, isolated tests. Each test gets a fresh database.

**Bug found — unchecked error** (`sqlite.go:225-227`): `CreateEvent()` fires an UPDATE to change parcel status but ignores the error return:
```go
s.db.ExecContext(ctx,
    "UPDATE parcels SET status = ?, updated_at = ? WHERE id = ?",
    e.Status, time.Now().UTC(), e.ParcelID)
// Error not checked!
```

### 3. `internal/handler` — High Priority (API Correctness)

**Why:** Contains business logic beyond routing — deduplication, background refresh, validation.

**Tests to add:**
- `RefreshParcel()` — event deduplication via `EventKey()`, re-fetches parcel mid-operation to avoid stale status
- `CreateParcel()` — validation (tracking_number required), defaults carrier to "manual", background refresh goroutine
- `Setup()` — calls `auth.Setup` then `auth.Login` in sequence, handles partial failure
- `CheckAuth()` — distinguishes "not configured" (needs setup) from "unauthorized"
- All endpoints — 400 (bad input), 404 (missing resource), 500 (store error) paths
- `writeJSON()`, `writeError()`, `decodeJSON()` — utility functions (easy wins)

**Approach:** Mock the `Store` interface and `tracker.Registry`. Test handlers via `httptest.NewRecorder`.

**Bug found — silent goroutine errors** (`parcels.go:99`): `CreateParcel` spawns a background goroutine for auto-refresh that silently discards errors — tracking failures are invisible.

### 4. `internal/tracker` — Worker & Registry (Medium Priority)

**Why:** The worker runs continuously in production; deduplication bugs mean duplicate events.

**Tests to add:**
- `Registry` — `Register()`, `Get()`, `Available()` basic operations
- `NewRegistry()` — all expected carriers are registered
- `EventKey()` — generates consistent dedup keys from timestamp+status+message
- `Worker.refreshAll()` — skips archived parcels, terminal statuses, manual carrier
- `Worker.refreshParcel()` — deduplicates events, updates last_check
- All carrier `Track()` methods — end-to-end with httptest mocking (currently 0%)

**Approach:** Mock Store + mock Tracker implementations. For the worker, use short intervals and context cancellation.

### 5. `internal/config` — Low Priority (Simple but Useful)

**Why:** Silent fallback on parse errors could cause surprising production behavior.

**Tests to add:**
- `Load()` with all env vars set, with none set (defaults), with invalid values
- `envInt()` / `envDuration()` — verify they silently fall back on parse errors (document this behavior)

**Approach:** Set `os.Setenv` in tests, defer cleanup with `t.Cleanup`.

### 6. Frontend — Low Priority (Requires Infrastructure Setup)

**Why:** No test framework exists. Setup cost is high but enables testing complex components.

**Infrastructure needed:** Vitest + `@testing-library/svelte` + `jsdom`

**Tests to add (once infrastructure exists):**
- `lib/utils.ts` — `formatRelativeTime()` edge cases (just now, minutes, hours, days)
- `lib/api.ts` — error handling, 401 redirect logic, type parsing
- `lib/i18n.svelte.ts` — locale detection, fallbacks, parameter substitution
- `pages/ParcelDetail.svelte` — most complex component (edit, refresh, archive, delete)

---

## Summary

The most impactful improvements, in order:
1. **Auth tests** — prevent security regressions, fix the race condition
2. **Store tests** — prevent data bugs, fix the unchecked error in CreateEvent
3. **Handler tests** — cover business logic and error paths
4. **Worker/Registry tests** — cover background processing and deduplication
5. **Config tests** — document silent fallback behavior
6. **Frontend test setup** — long-term investment for UI reliability
