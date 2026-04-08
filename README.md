# Parcel Tracker

A self-hosted parcel tracking application. Track all your packages in one place with a clean, modern interface.

## Features

- Track parcels from multiple carriers in a single dashboard
- Automatic background tracking — refreshes at configurable intervals
- Manual tracking events for any carrier
- Parcel archiving and search/filter support
- Clean, responsive UI with dark/light theme
- Single-binary deployment — no external dependencies
- SQLite database — zero configuration, easy backups
- Docker support — deploy anywhere in seconds
- Password-protected — simple single-user authentication

## Supported Carriers

| Carrier | API Key Required | Notes |
|---|---|---|
| La Poste | Yes (`LAPOSTE_API_KEY`) | French postal service |
| Colissimo | Yes (`LAPOSTE_API_KEY`) | Uses La Poste API |
| Chronopost | No | SOAP API |
| Mondial Relay | No | Tracking format: `expeditionNumber-postalCode` |
| GLS | No | Public REST API |
| DPD | No | |
| Colis Privé | No | |
| Relais Colis | No | |
| Manual | No | Add tracking events manually for any carrier |

## Quick Start

### Docker (recommended)

```bash
docker compose up -d
```

Then open [http://localhost:8080](http://localhost:8080) and set your password.

### Docker with environment password

```bash
docker compose up -d
# or
docker run -d -p 8080:8080 -v parcel-data:/data -e PARCEL_TRACKER_PASSWORD=yourpassword parcel-tracker
```

### Pre-built image

A pre-built image is published to `ghcr.io/simonbrunou/parcel-tracker` on every push to main and on release tags.

```bash
docker compose -f docker-compose.prod.yml up -d
```

### From source

Requirements: Go 1.25+, Node.js 22+

```bash
# Build everything
make build

# Run
./bin/parcel-tracker
```

## Configuration

| Environment Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `data/parcel-tracker.db` | SQLite database file path |
| `PARCEL_TRACKER_PASSWORD` | _(none)_ | Set initial password on first run |
| `REFRESH_INTERVAL` | `30m` | Background tracking refresh interval (set to `0` to disable) |
| `LAPOSTE_API_KEY` | _(none)_ | API key for La Poste/Colissimo tracking ([get one here](https://developer.laposte.fr)) |

## Development

```bash
# Install frontend dependencies
cd web && npm install && cd ..

# Run both frontend (Vite HMR) and backend concurrently
make dev

# Run tests
make test
```

- Frontend dev server: http://localhost:5173 (with HMR)
- Backend API: http://localhost:8080

## Architecture

```
cmd/parcel-tracker/     Go entrypoint
internal/
  auth/                 JWT authentication
  config/               Environment configuration
  handler/              HTTP handlers (REST API)
  model/                Data models
  server/               HTTP server + SPA handler
  store/                SQLite data layer
  tracker/              Carrier tracking implementations
web/
  src/                  Svelte 5 frontend
  dist/                 Built frontend (embedded in Go binary)
```

**Stack**: Go 1.25, Svelte 5, Tailwind CSS 4, Vite 6, TypeScript, SQLite

## API

All endpoints under `/api` require authentication (JWT cookie) except auth and health endpoints.

| Method | Path | Description |
|---|---|---|
| GET | `/api/health` | Health check + available carriers |
| POST | `/api/auth/setup` | Set initial password |
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/logout` | Logout |
| GET | `/api/auth/check` | Check auth status |
| GET | `/api/parcels` | List parcels (supports `status`, `search`, `archived` query params) |
| POST | `/api/parcels` | Create parcel |
| GET | `/api/parcels/:id` | Get parcel |
| PUT | `/api/parcels/:id` | Update parcel |
| DELETE | `/api/parcels/:id` | Delete parcel |
| POST | `/api/parcels/:id/refresh` | Refresh tracking |
| GET | `/api/parcels/:id/events` | List tracking events |
| POST | `/api/parcels/:id/events` | Add tracking event |
| DELETE | `/api/parcels/:id/events/:eventID` | Delete event |

## License

GPL-3.0
