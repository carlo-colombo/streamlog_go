package main

import (
	"bufio"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"log"
	"maps"
	"slices"

	"io"
)

type InMemoryLogsStore struct {
	logs    []logentry.Log
	logsCh  chan logentry.Log
	clients map[string]chan logentry.Log
}

func NewStore() *InMemoryLogsStore {
	return &InMemoryLogsStore{
		logsCh:  make(chan logentry.Log),
		logs:    []logentry.Log{},
		clients: make(map[string]chan logentry.Log),
	}
}

func (s *InMemoryLogsStore) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		log := logentry.Log{Line: line}
		s.logs = append(s.logs, log)
		fmt.Println(s.Clients())
		for _, client := range s.clients {
			client <- logentry.Log{Line: line}
		}
	}
}

func (s *InMemoryLogsStore) List() []logentry.Log {
	return s.logs
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

type Store interface {
	Scan(r io.Reader)
	List() []logentry.Log
	Disconnect(uid string)
	LineFor(uid string) chan logentry.Log
	Clients() []string
}
