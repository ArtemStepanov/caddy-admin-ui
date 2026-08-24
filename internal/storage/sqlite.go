package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage implements storage using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage
func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path+"?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	s := &SQLiteStorage{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if path != ":memory:" {
		if err := os.Chmod(path, 0600); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("secure database permissions: %w", err)
		}
	}

	return s, nil
}

func (s *SQLiteStorage) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS routes (
			id TEXT PRIMARY KEY,
			domain TEXT NOT NULL,
			path TEXT DEFAULT '',
			handler_type TEXT NOT NULL,
			config TEXT NOT NULL,
			enabled INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS global_config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS snapshots (
			id TEXT PRIMARY KEY,
			server TEXT NOT NULL,
			etag TEXT NOT NULL,
			routes TEXT NOT NULL,
			reason TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_routes_domain ON routes(domain);
		CREATE INDEX IF NOT EXISTS idx_routes_enabled ON routes(enabled);
	`)
	if err != nil {
		return err
	}

	columns := []struct {
		name       string
		definition string
	}{
		{"raw_caddy_route", "TEXT"},
		{"strip_path_prefix", "TEXT DEFAULT ''"},
		{"support_status", "TEXT DEFAULT ''"},
		{"readonly_reason", "TEXT DEFAULT ''"},
		{"position", "INTEGER DEFAULT 0"},
	}
	for _, column := range columns {
		exists, err := s.hasColumn("routes", column.name)
		if err != nil {
			return err
		}
		if !exists {
			if _, err := s.db.Exec(fmt.Sprintf("ALTER TABLE routes ADD COLUMN %s %s", column.name, column.definition)); err != nil {
				return fmt.Errorf("add routes.%s: %w", column.name, err)
			}
		}
	}
	_, err = s.db.Exec(`UPDATE routes
		SET support_status = ?, readonly_reason = CASE WHEN COALESCE(readonly_reason, '') = '' THEN ? ELSE readonly_reason END
		WHERE COALESCE(raw_caddy_route, '') != '' AND COALESCE(support_status, '') = ''`, SupportStatusPartialReadOnly, legacyReadOnlyReason)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) hasColumn(table, column string) (bool, error) {
	rows, err := s.db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}

// Route CRUD operations

// CreateRoute creates a new route
func (s *SQLiteStorage) CreateRoute(route *Route) error {
	return s.insertRoute(s.db, route)
}

// GetRoute retrieves a route by ID
func (s *SQLiteStorage) GetRoute(id string) (*Route, error) {
	row := s.db.QueryRow(
		`SELECT id, domain, path, handler_type, config, enabled, created_at, updated_at, COALESCE(raw_caddy_route, ''), COALESCE(strip_path_prefix, ''), COALESCE(support_status, ''), COALESCE(readonly_reason, ''), COALESCE(position, 0)
		 FROM routes WHERE id = ?`, id,
	)
	return s.scanRoute(row)
}

// ListRoutes returns all routes
func (s *SQLiteStorage) ListRoutes() ([]*Route, error) {
	rows, err := s.db.Query(
		`SELECT id, domain, path, handler_type, config, enabled, created_at, updated_at, COALESCE(raw_caddy_route, ''), COALESCE(strip_path_prefix, ''), COALESCE(support_status, ''), COALESCE(readonly_reason, ''), COALESCE(position, 0)
		 FROM routes ORDER BY position, created_at, id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*Route
	for rows.Next() {
		route, err := s.scanRouteRows(rows)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

// UpdateRoute updates an existing route
func (s *SQLiteStorage) UpdateRoute(route *Route) error {
	route.UpdatedAt = time.Now()
	_, err := s.db.Exec(
		`UPDATE routes SET domain=?, path=?, handler_type=?, config=?, enabled=?, updated_at=?, raw_caddy_route=?, strip_path_prefix=?, support_status=?, readonly_reason=?, position=?
		 WHERE id=?`,
		route.Domain, route.Path, route.HandlerType,
		string(route.Config), boolToInt(route.Enabled), route.UpdatedAt, string(route.RawCaddyRoute), route.StripPathPrefix, route.SupportStatus, route.ReadOnlyReason, route.Position, route.ID,
	)
	return err
}

// DeleteRoute deletes a route
func (s *SQLiteStorage) DeleteRoute(id string) error {
	_, err := s.db.Exec(`DELETE FROM routes WHERE id=?`, id)
	return err
}

// ReplaceAllRoutes atomically replaces all routes.
func (s *SQLiteStorage) ReplaceAllRoutes(routes []*Route) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM routes`); err != nil {
		return err
	}
	for _, route := range routes {
		if err := s.insertRoute(tx, route); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteStorage) insertRoute(exec interface {
	Exec(query string, args ...any) (sql.Result, error)
}, route *Route) error {
	if route.ID == "" {
		route.ID = uuid.New().String()
	}
	if route.SupportStatus == "" {
		route.SupportStatus = SupportStatusEditable
	}
	if route.CreatedAt.IsZero() {
		route.CreatedAt = time.Now()
	}
	route.UpdatedAt = time.Now()

	_, err := exec.Exec(
		`INSERT INTO routes (id, domain, path, handler_type, config, enabled, created_at, updated_at, raw_caddy_route, strip_path_prefix, support_status, readonly_reason, position)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		route.ID, route.Domain, route.Path, route.HandlerType,
		string(route.Config), boolToInt(route.Enabled), route.CreatedAt, route.UpdatedAt,
		string(route.RawCaddyRoute), route.StripPathPrefix, route.SupportStatus, route.ReadOnlyReason, route.Position,
	)
	return err
}

func (s *SQLiteStorage) scanRoute(row *sql.Row) (*Route, error) {
	var route Route
	var config string
	var enabled int
	var rawCaddyRoute string
	var stripPathPrefix string
	err := row.Scan(
		&route.ID, &route.Domain, &route.Path, &route.HandlerType,
		&config, &enabled, &route.CreatedAt, &route.UpdatedAt,
		&rawCaddyRoute, &stripPathPrefix, &route.SupportStatus, &route.ReadOnlyReason, &route.Position,
	)
	if err != nil {
		return nil, err
	}
	finishRouteScan(&route, config, enabled, rawCaddyRoute, stripPathPrefix)
	return &route, nil
}

func (s *SQLiteStorage) scanRouteRows(rows *sql.Rows) (*Route, error) {
	var route Route
	var config string
	var enabled int
	var rawCaddyRoute string
	var stripPathPrefix string
	err := rows.Scan(
		&route.ID, &route.Domain, &route.Path, &route.HandlerType,
		&config, &enabled, &route.CreatedAt, &route.UpdatedAt,
		&rawCaddyRoute, &stripPathPrefix, &route.SupportStatus, &route.ReadOnlyReason, &route.Position,
	)
	if err != nil {
		return nil, err
	}
	finishRouteScan(&route, config, enabled, rawCaddyRoute, stripPathPrefix)
	return &route, nil
}

func finishRouteScan(route *Route, config string, enabled int, rawCaddyRoute, stripPathPrefix string) {
	route.Config = json.RawMessage(config)
	route.Enabled = enabled == 1
	if rawCaddyRoute != "" {
		route.RawCaddyRoute = json.RawMessage(rawCaddyRoute)
	}
	route.StripPathPrefix = stripPathPrefix
	if route.SupportStatus == "" {
		if len(route.RawCaddyRoute) > 0 {
			route.SupportStatus = SupportStatusPartialReadOnly
		} else {
			route.SupportStatus = SupportStatusEditable
		}
	}
	if len(route.RawCaddyRoute) > 0 && route.IsReadOnly() && route.ReadOnlyReason == "" {
		route.ReadOnlyReason = legacyReadOnlyReason
	}
	route.ReadOnly = route.IsReadOnly()
}

// Global config

// GetGlobalConfig retrieves the global configuration
func (s *SQLiteStorage) GetGlobalConfig() (*GlobalConfig, error) {
	row := s.db.QueryRow(`SELECT value FROM global_config WHERE key = 'main'`)
	var value string
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		// URL default is owned by the API handler so env config is honored.
		return &GlobalConfig{EnableEncode: true}, nil
	}
	if err != nil {
		return nil, err
	}

	var config GlobalConfig
	if err := json.Unmarshal([]byte(value), &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SetGlobalConfig saves the global configuration
func (s *SQLiteStorage) SetGlobalConfig(config *GlobalConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT OR REPLACE INTO global_config (key, value) VALUES ('main', ?)`,
		string(data),
	)
	return err
}

// CreateSnapshot persists a pre-write copy of the managed route array.
func (s *SQLiteStorage) CreateSnapshot(snapshot *Snapshot) error {
	if snapshot.ID == "" {
		snapshot.ID = uuid.New().String()
	}
	if snapshot.CreatedAt.IsZero() {
		snapshot.CreatedAt = time.Now()
	}
	_, err := s.db.Exec(`INSERT INTO snapshots (id, server, etag, routes, reason, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		snapshot.ID, snapshot.Server, snapshot.ETag, string(snapshot.Routes), snapshot.Reason, snapshot.CreatedAt)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM snapshots WHERE id IN (SELECT id FROM snapshots ORDER BY created_at DESC LIMIT -1 OFFSET 100)`)
	return err
}

// ListSnapshots returns newest snapshots first without their potentially large route payload.
func (s *SQLiteStorage) ListSnapshots() ([]*Snapshot, error) {
	rows, err := s.db.Query(`SELECT id, server, etag, reason, created_at FROM snapshots ORDER BY created_at DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snapshots []*Snapshot
	for rows.Next() {
		snapshot := &Snapshot{}
		if err := rows.Scan(&snapshot.ID, &snapshot.Server, &snapshot.ETag, &snapshot.Reason, &snapshot.CreatedAt); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, rows.Err()
}

// GetSnapshot returns a snapshot including its route payload.
func (s *SQLiteStorage) GetSnapshot(id string) (*Snapshot, error) {
	snapshot := &Snapshot{}
	var routes string
	err := s.db.QueryRow(`SELECT id, server, etag, routes, reason, created_at FROM snapshots WHERE id = ?`, id).
		Scan(&snapshot.ID, &snapshot.Server, &snapshot.ETag, &routes, &snapshot.Reason, &snapshot.CreatedAt)
	if err != nil {
		return nil, err
	}
	snapshot.Routes = json.RawMessage(routes)
	return snapshot, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
