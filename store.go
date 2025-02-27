package main

import (
	"bufio"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"log"
	"maps"
	"slices"

	"io"
)

type Store struct {
	logs    []logentry.Log
	logsCh  chan logentry.Log
	clients map[string]chan logentry.Log
}

func NewStore() Store {
	return Store{
		logsCh:  make(chan logentry.Log),
		logs:    []logentry.Log{},
		clients: make(map[string]chan logentry.Log),
	}
}

func (s *Store) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		log := logentry.Log{Line: line}
		s.logs = append(s.logs, log)
		for _, client := range s.clients {
			client <- logentry.Log{Line: line}
		}
	}
}

func (s *Store) List() []logentry.Log {
	return s.logs
}

func (s *Store) Disconnect(uid string) {
	delete(s.clients, uid)
	log.Printf("Client %s disconnected", uid)
}

func (s *Store) LineFor(uid string) chan logentry.Log {
	if _, ok := s.clients[uid]; !ok {
		s.clients[uid] = make(chan logentry.Log)
	}

	return s.clients[uid]
}

func (s *Store) Clients() []string {
	return slices.Sorted(maps.Keys(s.clients))
}
