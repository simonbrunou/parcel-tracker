# Test Coverage Analysis

_Updated: 2026-04-09_

## Coverage Summary

| Package | Before | After | Tests Added |
|---|---|---|---|
| `internal/auth` | 0% | **86.7%** | 15 tests |
| `internal/config` | 0% | **100.0%** | 4 tests |
| `internal/handler` | 0% | **72.1%** | 28 tests |
| `internal/server` | 0% | **97.6%** | 6 tests |
| `internal/store` | 0% | **86.6%** | 18 tests |
| `internal/tracker` | 53.1% | **65.8%** | 14 tests |
| `internal/model` | N/A | N/A | Constants only |
| `web/` (frontend) | 0% | 0% | No test framework |

## Bugs Fixed

### 1. Unchecked error in `CreateEvent` (`store/sqlite.go:225`)

`CreateEvent()` updated parcel status but ignored the error return. Now properly checked.

### 2. Race condition in `jwtSecret` (`auth/auth.go:131`)

Multiple goroutines calling `jwtSecret()` simultaneously could generate and save different secrets, invalidating tokens. Added `sync.Mutex` to serialize secret generation.

### 3. Foreign keys not enforced (`store/sqlite.go:27`)

The `_foreign_keys=1` connection parameter wasn't recognized by the `modernc.org/sqlite` driver. Changed to `_pragma=foreign_keys(1)` format which is properly applied. This means `ON DELETE CASCADE` on tracking events now works correctly.

### 4. SPA cache header never set (`server/spa.go:26`)

The check `strings.Contains(path, "/assets/")` never matched because the leading slash was already trimmed. Changed to `strings.HasPrefix(path, "assets/")`.

## Remaining Gaps

- **Tracker `Track()` methods**: End-to-end HTTP calls for each carrier remain untested
- **Handler background goroutine**: `CreateParcel` auto-refresh silently discards errors
- **Frontend**: No test framework (would need Vitest + @testing-library/svelte)
