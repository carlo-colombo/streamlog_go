package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	stdlog "log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/carlo-colombo/streamlog_go/logentry"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteLogsStore struct {
	db             *sql.DB
	clients        map[string]chan logentry.Log
	filter         string
	filterChangeCh chan struct{}
}

func retryWithBackoff(operation func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		backoff := time.Duration(1<<uint(i)) * time.Second
		stdlog.Printf("Operation failed, retrying in %v: %v", backoff, err)
		time.Sleep(backoff)
	}
	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, err)
}

func NewSQLiteStore(dbPath string) (*SQLiteLogsStore, error) {
	var db *sql.DB
	err := retryWithBackoff(func() error {
		var err error
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		return nil
	}, 3)

	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Create logs table if it doesn't exist with transaction
	err = retryWithBackoff(func() error {
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback()

		// Create table if it doesn't exist
		_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS logs (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				line TEXT NOT NULL,
				timestamp DATETIME NOT NULL
			)
		`)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		// Verify table exists and is accessible
		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to verify table: %w", err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}, 3)

	if err != nil {
		return nil, fmt.Errorf("table creation failed: %w", err)
	}

	return &SQLiteLogsStore{
		db:             db,
		clients:        make(map[string]chan logentry.Log),
		filterChangeCh: make(chan struct{}),
	}, nil
}

func (s *SQLiteLogsStore) SetFilter(filter string) {
	s.filter = filter
	if len(s.clients) > 0 {
		s.filterChangeCh <- struct{}{}
	}
}

func (s *SQLiteLogsStore) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		logLine := logentry.NewLog(line)

		// Insert log into database with retry
		err := retryWithBackoff(func() error {
			_, err := s.db.Exec(
				"INSERT INTO logs (line, timestamp) VALUES (?, ?)",
				logLine.Line,
				logLine.Timestamp,
			)
			if err != nil {
				return fmt.Errorf("failed to insert log: %w", err)
			}
			return nil
		}, 3)

		if err != nil {
			stdlog.Printf("Failed to insert log after retries: %v", err)
			continue
		}

		// Broadcast to clients if it matches the filter
		if s.filter == "" || strings.Contains(strings.ToLower(logLine.Line), strings.ToLower(s.filter)) {
			for _, client := range s.clients {
				client <- logLine
			}
		}
	}
}

func (s *SQLiteLogsStore) List() []logentry.Log {
	var query string
	var args []interface{}

	if s.filter != "" {
		// Use REPLACE to add ANSI highlighting to matched terms
		query = `
			SELECT 
				REPLACE(
					REPLACE(
						line,
						LOWER(?),
						CHAR(27) || '[43m' || LOWER(?) || CHAR(27) || '[0m'
					),
					UPPER(?),
					CHAR(27) || '[43m' || UPPER(?) || CHAR(27) || '[0m'
				) as line,
				timestamp 
			FROM logs 
			WHERE LOWER(line) LIKE LOWER(?) 
			ORDER BY id ASC`
		args = []interface{}{s.filter, s.filter, s.filter, s.filter, "%" + s.filter + "%"}
	} else {
		query = "SELECT line, timestamp FROM logs ORDER BY id ASC"
	}

	var logs []logentry.Log
	err := retryWithBackoff(func() error {
		rows, err := s.db.Query(query, args...)
		if err != nil {
			return fmt.Errorf("failed to query logs: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var log logentry.Log
			err := rows.Scan(&log.Line, &log.Timestamp)
			if err != nil {
				return fmt.Errorf("failed to scan log: %w", err)
			}
			logs = append(logs, log)
		}
		return rows.Err()
	}, 10)

	if err != nil {
		stdlog.Printf("Failed to list logs after retries: %v", err)
		return nil
	}

	return logs
}

func (s *SQLiteLogsStore) Disconnect(uid string) {
	delete(s.clients, uid)
	stdlog.Printf("Client %s disconnected", uid)
}

func (s *SQLiteLogsStore) LineFor(uid string) chan logentry.Log {
	if _, ok := s.clients[uid]; !ok {
		s.clients[uid] = make(chan logentry.Log)
	}
	return s.clients[uid]
}

func (s *SQLiteLogsStore) Clients() []string {
	return slices.Sorted(maps.Keys(s.clients))
}

func (s *SQLiteLogsStore) FilterChangeFor() chan struct{} {
	return s.filterChangeCh
}

func (s *SQLiteLogsStore) Close() error {
	return s.db.Close()
}

type Store interface {
	SetFilter(filter string)
	Scan(r io.Reader)
	List() []logentry.Log
	Disconnect(uid string)
	LineFor(uid string) chan logentry.Log
	Clients() []string
	FilterChangeFor() chan struct{}
}
