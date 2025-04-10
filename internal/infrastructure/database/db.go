package database

import (
	"context"
	"time"

	"easy-go-backend/ent"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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

		log.Warn().Int("attempt", i+1).Int("maxRetries", maxRetries).Err(err).Msg("Failed to connect to database")
		time.Sleep(time.Duration(i+1) * time.Second)
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
