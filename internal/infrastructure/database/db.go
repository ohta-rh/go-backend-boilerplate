package database

import (
	"context"
	"log"
	"time"

	"easy-go-backend/ent"

	_ "github.com/lib/pq"
)

// NewClient creates a new ent client with connection to PostgreSQL.
func NewClient(dsn string, maxRetries int) (*ent.Client, error) {
	var client *ent.Client

	var err error

	for i := range maxRetries {
		client, err = ent.Open("postgres", dsn)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, err
	}

	// Run schema migration
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}
