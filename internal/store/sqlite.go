package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/simonbrunou/parcel-tracker/internal/model"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return s, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS parcels (
			id TEXT PRIMARY KEY,
			tracking_number TEXT NOT NULL,
			carrier TEXT NOT NULL DEFAULT 'manual',
			name TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'unknown',
			archived INTEGER NOT NULL DEFAULT 0,
			last_check DATETIME,
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);

		CREATE TABLE IF NOT EXISTS tracking_events (
			id TEXT PRIMARY KEY,
			parcel_id TEXT NOT NULL REFERENCES parcels(id) ON DELETE CASCADE,
			status TEXT NOT NULL DEFAULT 'unknown',
			message TEXT NOT NULL DEFAULT '',
			location TEXT NOT NULL DEFAULT '',
			timestamp DATETIME NOT NULL DEFAULT (datetime('now')),
			created_at DATETIME NOT NULL DEFAULT (datetime('now'))
		);

		CREATE INDEX IF NOT EXISTS idx_tracking_events_parcel_id ON tracking_events(parcel_id);
		CREATE INDEX IF NOT EXISTS idx_parcels_status ON parcels(status);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_parcels_tracking_carrier ON parcels(tracking_number, carrier) WHERE archived = 0;
		CREATE INDEX IF NOT EXISTS idx_parcels_archived ON parcels(archived);
		CREATE INDEX IF NOT EXISTS idx_parcels_updated_at ON parcels(updated_at);

		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL DEFAULT ''
		);
	`)
	return err
}

func newID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}

// Parcels

func (s *SQLiteStore) ListParcels(ctx context.Context, filter ParcelFilter) (PaginatedParcels, error) {
	where := " WHERE 1=1"
	var args []any

	if filter.Status != "" {
		where += " AND status = ?"
		args = append(args, filter.Status)
	}
	if filter.Archived != nil {
		where += " AND archived = ?"
		if *filter.Archived {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if filter.Search != "" {
		where += " AND (name LIKE ? OR tracking_number LIKE ?)"
		searchPat := "%" + filter.Search + "%"
		args = append(args, searchPat, searchPat)
	}

	// Count total matching rows.
	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM parcels"+where, args...).Scan(&total); err != nil {
		return PaginatedParcels{}, err
	}

	// Apply pagination defaults.
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	query := "SELECT id, tracking_number, carrier, name, notes, status, archived, last_check, created_at, updated_at FROM parcels" + where + " ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	paginatedArgs := append(args, pageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, paginatedArgs...)
	if err != nil {
		return PaginatedParcels{}, err
	}
	defer rows.Close()

	var parcels []model.Parcel
	for rows.Next() {
		p, err := scanParcel(rows)
		if err != nil {
			return PaginatedParcels{}, err
		}
		parcels = append(parcels, p)
	}
	if parcels == nil {
		parcels = []model.Parcel{}
	}
	if err := rows.Err(); err != nil {
		return PaginatedParcels{}, err
	}

	return PaginatedParcels{
		Data:     parcels,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *SQLiteStore) ListActiveParcels(ctx context.Context) ([]model.Parcel, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, tracking_number, carrier, name, notes, status, archived, last_check, created_at, updated_at
		 FROM parcels
		 WHERE archived = 0 AND status NOT IN (?, ?) AND carrier != ?
		 ORDER BY updated_at DESC`,
		model.StatusDelivered, model.StatusExpired, model.CarrierManual)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []model.Parcel
	for rows.Next() {
		p, err := scanParcel(rows)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	if parcels == nil {
		parcels = []model.Parcel{}
	}
	return parcels, rows.Err()
}

func (s *SQLiteStore) GetParcel(ctx context.Context, id string) (model.Parcel, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, tracking_number, carrier, name, notes, status, archived, last_check, created_at, updated_at FROM parcels WHERE id = ?", id)
	return scanParcelRow(row)
}

func (s *SQLiteStore) CreateParcel(ctx context.Context, p model.Parcel) (model.Parcel, error) {
	p.ID = newID()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Status == "" {
		p.Status = model.StatusUnknown
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO parcels (id, tracking_number, carrier, name, notes, status, archived, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.TrackingNumber, p.Carrier, p.Name, p.Notes, p.Status, p.Archived, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return model.Parcel{}, err
	}
	return p, nil
}

func (s *SQLiteStore) UpdateParcel(ctx context.Context, p model.Parcel) (model.Parcel, error) {
	p.UpdatedAt = time.Now().UTC()
	_, err := s.db.ExecContext(ctx,
		`UPDATE parcels SET tracking_number = ?, carrier = ?, name = ?, notes = ?, status = ?, archived = ?, last_check = ?, updated_at = ?
		WHERE id = ?`,
		p.TrackingNumber, p.Carrier, p.Name, p.Notes, p.Status, p.Archived, p.LastCheck, p.UpdatedAt, p.ID)
	if err != nil {
		return model.Parcel{}, err
	}
	return p, nil
}

func (s *SQLiteStore) DeleteParcel(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM parcels WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Tracking Events

func (s *SQLiteStore) ListEvents(ctx context.Context, parcelID string) ([]model.TrackingEvent, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, parcel_id, status, message, location, timestamp, created_at FROM tracking_events WHERE parcel_id = ? ORDER BY timestamp DESC",
		parcelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.TrackingEvent
	for rows.Next() {
		var e model.TrackingEvent
		if err := rows.Scan(&e.ID, &e.ParcelID, &e.Status, &e.Message, &e.Location, &e.Timestamp, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if events == nil {
		events = []model.TrackingEvent{}
	}
	return events, rows.Err()
}

func (s *SQLiteStore) CreateEvent(ctx context.Context, e model.TrackingEvent) (model.TrackingEvent, error) {
	e.ID = newID()
	e.CreatedAt = time.Now().UTC()
	if e.Timestamp.IsZero() {
		e.Timestamp = e.CreatedAt
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO tracking_events (id, parcel_id, status, message, location, timestamp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.ParcelID, e.Status, e.Message, e.Location, e.Timestamp, e.CreatedAt)
	if err != nil {
		return model.TrackingEvent{}, err
	}

	// Update parcel status and timestamp
	if _, err := s.db.ExecContext(ctx,
		"UPDATE parcels SET status = ?, updated_at = ? WHERE id = ?",
		e.Status, time.Now().UTC(), e.ParcelID); err != nil {
		return e, fmt.Errorf("update parcel status: %w", err)
	}

	return e, nil
}

func (s *SQLiteStore) DeleteEvent(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM tracking_events WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Settings

func (s *SQLiteStore) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, "SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *SQLiteStore) SetSetting(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value)
	return err
}

// Helpers

type scanner interface {
	Scan(dest ...any) error
}

func scanParcelFields(s scanner) (model.Parcel, error) {
	var p model.Parcel
	var lastCheck sql.NullTime
	var archived int
	err := s.Scan(&p.ID, &p.TrackingNumber, &p.Carrier, &p.Name, &p.Notes, &p.Status, &archived, &lastCheck, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return model.Parcel{}, err
	}
	p.Archived = archived != 0
	if lastCheck.Valid {
		p.LastCheck = &lastCheck.Time
	}
	return p, nil
}

func scanParcel(rows *sql.Rows) (model.Parcel, error) {
	return scanParcelFields(rows)
}

func scanParcelRow(row *sql.Row) (model.Parcel, error) {
	return scanParcelFields(row)
}

// Helper to filter out empty strings in query building
func buildSearchPattern(search string) string {
	return "%" + strings.ReplaceAll(search, "%", "\\%") + "%"
}
