package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/carlo-colombo/streamlog_go/logentry"

	"io"
)

type InMemoryLogsStore struct {
	logs           []logentry.Log
	logsCh         chan logentry.Log
	clients        map[string]chan logentry.Log
	filter         string
	filterChangeCh chan struct{}
	uid            string
}

func NewStore() *InMemoryLogsStore {
	return &InMemoryLogsStore{
		logsCh:         make(chan logentry.Log),
		logs:           []logentry.Log{},
		clients:        make(map[string]chan logentry.Log),
		filter:         "",
		filterChangeCh: make(chan struct{}),
		uid:            fmt.Sprintf("client-%d", time.Now().UnixNano()),
	}
}

func (s *InMemoryLogsStore) SetFilter(filter string) {
	s.filter = filter
	if len(s.clients) > 0 {
		s.filterChangeCh <- struct{}{}
	}
}

func (s *InMemoryLogsStore) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		logLine := logentry.NewLog(line)
		s.logs = append(s.logs, logLine)
		if s.filter == "" || strings.Contains(strings.ToLower(logLine.Line), strings.ToLower(s.filter)) {
			for _, client := range s.clients {
				client <- logLine
			}
		}
	}
}

func (s *InMemoryLogsStore) List() []logentry.Log {
	if s.filter == "" {
		return s.logs
	}

	filtered := make([]logentry.Log, 0)
	for _, log := range s.logs {
		if strings.Contains(strings.ToLower(log.Line), strings.ToLower(s.filter)) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

func (s *InMemoryLogsStore) Disconnect(uid string) {
	delete(s.clients, uid)
	log.Printf("Client %s disconnected", uid)
}

func (s *InMemoryLogsStore) LineFor(uid string) chan logentry.Log {
	if _, ok := s.clients[uid]; !ok {
		s.clients[uid] = make(chan logentry.Log)
	}

	return s.clients[uid]
}

func (s *InMemoryLogsStore) Clients() []string {
	return slices.Sorted(maps.Keys(s.clients))
}

func (s *InMemoryLogsStore) FilterChangeFor() chan struct{} {
	return s.filterChangeCh
}

type Store interface {
	Scan(r io.Reader)
	List() []logentry.Log
	Disconnect(uid string)
	LineFor(uid string) chan logentry.Log
	Clients() []string
	SetFilter(filter string)
	FilterChangeFor() chan struct{}
}
