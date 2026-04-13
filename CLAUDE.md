# CLAUDE.md

## Project Overview

Parcel Tracker is a self-hosted, single-user parcel tracking application. It consolidates package tracking from multiple carriers into a dashboard with automatic background refresh. Built with a Go backend and Svelte 5 frontend, it deploys as a single binary with an embedded SQLite database.

## Tech Stack

- **Backend**: Go 1.25, chi router, SQLite (modernc.org/sqlite), JWT auth
- **Frontend**: Svelte 5, TypeScript, Tailwind CSS 4, Vite 6
- **Database**: SQLite with WAL mode, raw SQL (no ORM)
- **Deployment**: Docker multi-stage build, single binary with embedded frontend

## Repository Structure

```
cmd/parcel-tracker/          # Go entrypoint (main.go)
internal/
  auth/                      # JWT + bcrypt authentication
  config/                    # Environment variable configuration
  handler/                   # HTTP request handlers (REST API)
  model/                     # Data models (Parcel, TrackingEvent, enums)
  server/                    # HTTP server setup, chi router, SPA handler
  store/                     # SQLite data access layer (repository pattern)
  tracker/                   # Carrier tracking implementations + background worker
web/                         # Svelte 5 frontend
  src/
    components/              # Reusable UI (ParcelCard, StatusBadge, Navbar, ParcelTimeline)
    pages/                   # Route pages (Dashboard, Login, AddParcel, ParcelDetail)
    lib/                     # Shared utilities (api.ts, i18n.svelte.ts, utils.ts)
    App.svelte               # Root SPA router
  dist/                      # Built frontend (embedded into Go binary via embed.go)
.github/workflows/           # CI/CD (Docker image publish to GHCR)
```

## Common Commands

```bash
make dev              # Run backend + frontend dev servers concurrently
make dev-backend      # Go server on :8080 with DEV=1
make dev-frontend     # Vite HMR dev server on :5173 (proxies /api to :8080)
make build            # Full production build (frontend then backend)
make test             # Run all Go tests (go test ./...)
make docker           # Build Docker image
make clean            # Remove bin/, web/dist/, data/
```

Frontend commands (run from `web/` directory):
```bash
npm run dev           # Vite dev server
npm run build         # Production build to dist/
```

## Architecture

### Layered Design

1. **Handler layer** (`internal/handler/`): HTTP handlers receive requests, validate input, call store/tracker, return JSON
2. **Store layer** (`internal/store/`): `Store` interface abstracts SQLite operations. Implementation in `sqlite.go` with auto-migration
3. **Tracker layer** (`internal/tracker/`): `Tracker` interface per carrier. `Registry` holds all implementations. `Worker` runs background refresh
4. **Auth layer** (`internal/auth/`): Password hashing (bcrypt), JWT tokens (30-day lifetime), HttpOnly session cookies

### API Routes

All API routes are under `/api`. Auth routes (`/auth/*`) and `/health` are public. Everything under `/parcels` requires JWT authentication via middleware.

```
POST   /api/auth/setup           # Initial password setup
POST   /api/auth/login            # Login
POST   /api/auth/logout           # Logout
GET    /api/auth/check            # Auth status check
GET    /api/health                # Health + available carriers

# Protected (JWT required):
GET    /api/parcels               # List parcels (?status, ?search, ?archived)
POST   /api/parcels               # Create parcel
GET    /api/parcels/{id}          # Get parcel
PUT    /api/parcels/{id}          # Update parcel
DELETE /api/parcels/{id}          # Delete parcel
POST   /api/parcels/{id}/refresh  # Refresh tracking
GET    /api/parcels/{id}/events   # List events
POST   /api/parcels/{id}/events   # Create manual event
DELETE /api/parcels/{id}/events/{eventID}  # Delete event
```

Non-API routes serve the embedded Svelte SPA with index.html fallback.

### Database

SQLite with WAL mode, foreign keys enabled, 5s busy timeout. Tables: `parcels`, `tracking_events`, `settings`. Schema auto-migrates on startup. IDs are UUID v7. Raw SQL with `database/sql` prepared statements.

### Frontend

- **Routing**: `svelte-spa-router` with hash-based URLs (`#/login`, `#/parcels/:id`)
- **State**: Svelte 5 runes (`$state`, `$derived`, `$effect`) — no stores/context
- **API client**: `web/src/lib/api.ts` — typed fetch wrapper with auto 401 redirect
- **Theming**: Dark/light via `data-theme` attribute, stored in localStorage
- **i18n**: `web/src/lib/i18n.svelte.ts` with `domain.key` translation pattern, browser locale detection

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `data/parcel-tracker.db` | SQLite file path |
| `PARCEL_TRACKER_PASSWORD` | _(none)_ | Initial password (first run only) |
| `REFRESH_INTERVAL` | `30m` | Background tracking refresh interval |
| `LAPOSTE_API_KEY` | _(none)_ | La Poste/Colissimo API key |
| `DEV` | _(unset)_ | Dev mode flag |

## Adding a New Carrier

1. Create `internal/tracker/<carrier>.go` implementing the `Tracker` interface:
   - `Code() model.CarrierCode` — return the carrier code constant
   - `Name() string` — human-readable name
   - `Track(ctx, trackingNumber) ([]model.TrackingEvent, error)` — fetch and parse tracking data
2. Add a `CarrierCode` constant in `internal/model/parcel.go`
3. Register the tracker in `NewRegistry()` in `internal/tracker/tracker.go`
4. Add tests in `internal/tracker/<carrier>_test.go` using `httptest.NewServer` for HTTP mocking

## Testing Conventions

- Go standard `testing` package only — no external test framework
- HTTP mocking with `net/http/httptest`
- Tests exist for carrier tracker implementations (`internal/tracker/*_test.go`)
- No frontend tests currently
- Run: `make test`

## Code Conventions

### Go
- Standard Go formatting (`gofmt`)
- `log/slog` for structured logging
- Explicit error returns (`if err != nil`)
- All DB/tracker operations accept `context.Context`
- Handler request/response types as small structs in handler files
- JSON helper functions: `writeJSON()`, `writeError()`, `decodeJSON()`

### TypeScript / Svelte
- Strict TypeScript (`strict: true` in tsconfig)
- PascalCase for components, camelCase for variables/functions
- Svelte 5 runes for reactivity (not legacy stores)
- Tailwind utility classes for styling
- CSS custom properties for theme variables

## CI/CD

GitHub Actions (`.github/workflows/docker-publish.yml`):
- Triggers on pushes to `main`, version tags (`v*`), and manual dispatch
- Builds and pushes Docker image to `ghcr.io/simonbrunou/parcel-tracker`
- Uses Docker Buildx with GitHub Actions cache

## graphify

This project has a graphify knowledge graph at graphify-out/.

Rules:
- Before answering architecture or codebase questions, read graphify-out/GRAPH_REPORT.md for god nodes and community structure
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files
- After modifying code files in this session, run `python3 -c "from graphify.watch import _rebuild_code; from pathlib import Path; _rebuild_code(Path('.'))"` to keep the graph current
