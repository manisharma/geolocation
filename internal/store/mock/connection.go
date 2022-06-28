package mock

import (
	"context"
	"geolocation/internal/model"
	"geolocation/internal/utils"
)

// Connection mocks the phycial database connection with in-memory store
// implements internal.Store
type Connection struct {
	locations []model.Location
	closed    bool
}

// New creates new in-memory datastore
func New() (*Connection, error) {
	return &Connection{locations: []model.Location{}, closed: false}, nil
}

func (s *Connection) Migrate(ctx context.Context) error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (s *Connection) BulkCreate(ctx context.Context, locations []model.Location) error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.locations = locations
	return nil
}

func (s *Connection) Get(ctx context.Context, ipAddress string) (*model.Location, error) {
	if s.closed {
		return nil, utils.ErrInvalidConn
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for _, l := range s.locations {
		if l.IPAddress == ipAddress {
			return &l, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (s *Connection) Close() error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	s.locations = nil
	return nil
}
