package main

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func Test_AddEvent_Success(t *testing.T) {
	storage := NewMemoryStorage()

	event := Event{
		OrderType:  "Purchase",
		SessionID:  "29827525-06c9-4b1e-9d9b-7c4584e82f56",
		Card:       "4433**1409",
		EventDate:  "2023-01-04 13:44:52.835626 +00:00",
		WebsiteURL: "https://amazon.com",
	}

	err := storage.AddEvent(event)
	if err != nil {
		t.Fatalf("want nil got %v", err)
	}
}

func Test_GetEvents_Success(t *testing.T) {
	storage := NewMemoryStorage()

	event := Event{
		OrderType:  "Purchase",
		SessionID:  "29827525-06c9-4b1e-9d9b-7c4584e82f56",
		Card:       "4433**1409",
		EventDate:  "2023-01-04 13:44:52.835626 +00:00",
		WebsiteURL: "https://amazon.com",
	}

	err := storage.AddEvent(event)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	events, err := storage.GetEvents()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if !reflect.DeepEqual(events[0], event) {
		t.Errorf("expected %v event, got %v", events[0], event)
	}

}

func Test_GetEvents_Error(t *testing.T) {
	storage := NewMemoryStorage()

	_, err := storage.GetEvents()
	if !errors.Is(err, ErrEventsNotFound) {
		t.Errorf("expected error after reading all events, got %v", err)
	}
}

func Test_Events_Concurrency(t *testing.T) {
	storage := NewMemoryStorage()
	var wg sync.WaitGroup

	goroutineCount := 10
	eventCount := 100

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < eventCount; j++ {
				event := Event{
					OrderType:  "TestEvent",
					SessionID:  "session-" + strconv.Itoa(i) + "-" + strconv.Itoa(j),
					Card:       "4433**1409",
					EventDate:  "2023-01-04T13:44:52Z",
					WebsiteURL: "https://test.com",
				}
				if err := storage.AddEvent(event); err != nil {
					t.Errorf("AddEvent failed: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()

	events, err := storage.GetEvents()
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	expectedTotal := goroutineCount * eventCount
	if len(events) != expectedTotal {
		t.Errorf("Expected %d events, got %d", expectedTotal, len(events))
	}
}
