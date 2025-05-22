package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Job struct {
	storage Storage
}

func NewJob(storage Storage) *Job {
	return &Job{
		storage: storage,
	}
}

func (j *Job) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Job stopped")
			return
		case <-ticker.C:
			events, err := j.storage.GetEvents()
			if err != nil {
				log.Println("Job -> StartJob -> ", err)
			}
			for _, e := range events {
				fmt.Printf("Notify: %s - %s at %s\n", e.OrderType, e.Card, e.WebsiteURL)
			}
		}
	}
}
