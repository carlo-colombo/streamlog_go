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

	"github.com/carlo-colombo/streamlog_go/logentry"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteLogsStore struct {
	db             *sql.DB
	clients        map[string]chan logentry.Log
	filter         string
	filterChangeCh chan struct{}
}

func NewSQLiteStore(dbPath string) (*SQLiteLogsStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create logs table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			line TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
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

		// Insert log into database
		_, err := s.db.Exec(
			"INSERT INTO logs (line, timestamp) VALUES (?, ?)",
			logLine.Line,
			logLine.Timestamp,
		)
		if err != nil {
			stdlog.Printf("Failed to insert log: %v", err)
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
		query = "SELECT line, timestamp FROM logs WHERE LOWER(line) LIKE LOWER(?) ORDER BY timestamp DESC"
		args = []interface{}{"%" + s.filter + "%"}
	} else {
		query = "SELECT line, timestamp FROM logs ORDER BY timestamp DESC"
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		stdlog.Printf("Failed to query logs: %v", err)
		return nil
	}
	defer rows.Close()

	var logs []logentry.Log
	for rows.Next() {
		var log logentry.Log
		err := rows.Scan(&log.Line, &log.Timestamp)
		if err != nil {
			stdlog.Printf("Failed to scan log: %v", err)
			continue
		}
		logs = append(logs, log)
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
