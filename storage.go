package main

import (
	"errors"
	"sync"
)

type Storage interface {
	AddEvent(event Event) error
	GetEvents() ([]Event, error)
}

type MemoryStorage struct {
	events []Event
	mu     sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		events: make([]Event, 0),
	}
}

var ErrEventsNotFound = errors.New("events not found")

func (s *MemoryStorage) AddEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return nil
}

func (s *MemoryStorage) GetEvents() ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.events) == 0 {
		return nil, ErrEventsNotFound
	}

	events := s.events
	s.events = nil
	return events, nil
}
