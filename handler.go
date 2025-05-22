package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type EventHandler struct {
	storage Storage
}

func NewEventHandler(storage Storage) *EventHandler {
	return &EventHandler{storage: storage}
}

func (h *EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("EventHandler -> AddEvent -> Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("EventHandler -> AddEvent -> error %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	err := h.storage.AddEvent(event)
	if err != nil {
		log.Printf("EventHandler -> AddEvent -> error %s", err.Error())
		http.Error(w, "Error on add to storage", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
