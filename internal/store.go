package internal

import (
	"context"
	"geolocation/internal/model"
)

// Store is the common interface to query data from Database or In-memory store
type Store interface {
	// Migrate meant to run migration script (if any)
	Migrate(ctx context.Context) error
	// BulkCreate is meant to do bulk insertion of data
	BulkCreate(ctx context.Context, locations []model.Location) error
	// Get is meant to retrieve data on the basis of IP address
	Get(ctx context.Context, ipAddress string) (*model.Location, error)
	// Close is meant to close the connection
	Close() error
}
