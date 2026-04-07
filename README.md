# Parcel Tracker

A self-hosted parcel tracking application. Track all your packages in one place with a clean, modern interface.

## Features

- Track parcels from multiple carriers in a single dashboard
- Manual tracking events for any carrier
- Clean, responsive UI with dark/light theme
- Single-binary deployment — no external dependencies
- SQLite database — zero configuration, easy backups
- Docker support — deploy anywhere in seconds
- Password-protected — simple single-user authentication

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

### From source

Requirements: Go 1.24+, Node.js 22+

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

## Development

```bash
# Install frontend dependencies
cd web && npm install && cd ..

# Run both frontend (Vite HMR) and backend concurrently
make dev
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
  tracker/              Carrier tracking abstraction
web/
  src/                  Svelte 5 frontend
  dist/                 Built frontend (embedded in Go binary)
```

**Stack**: Go, Svelte 5, Tailwind CSS 4, SQLite

## API

All endpoints under `/api` require authentication (JWT cookie) except auth endpoints.

| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/setup` | Set initial password |
| POST | `/api/auth/login` | Login |
| POST | `/api/auth/logout` | Logout |
| GET | `/api/auth/check` | Check auth status |
| GET | `/api/parcels` | List parcels |
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
