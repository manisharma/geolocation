package database

import (
	"context"
	"errors"
	"fmt"
	"geolocation/internal/model"
	"geolocation/internal/utils"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Connection wraps a physical database connection
type Connection struct {
	db     *sqlx.DB
	closed bool
}

// New creates new databse connection
// implements internal.Store
func New(cfg model.DB) (*Connection, error) {
	connStr := "host=" + cfg.Host + " port=" + strconv.Itoa(cfg.Port) + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Database + " sslmode=disable binary_parameters=yes"
	db, err := sqlx.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Connection{db, false}, nil
}

// Migrate runs migration on database
func (s *Connection) Migrate(ctx context.Context) error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	_, err := s.db.Exec(utils.CreateTableScript)
	return err
}

// BulkCreate inserts locations data in database
func (s *Connection) BulkCreate(ctx context.Context, locations []model.Location) error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	var size int = 1000
	batches := [][]model.Location{}
	if size <= 0 || size > len(locations) {
		batches = append(batches, locations)
	} else {
		lenBulk := len(locations)
		for from := 0; from < lenBulk; from += size {
			till := from + size
			if till > lenBulk {
				till = lenBulk
			}
			batches = append(batches, locations[from:till])
		}
	}

	success := false
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, batch := range batches {
		vals := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*3)
		i := 0
		idx := 0
		for _, row := range batch {
			idx = i * 7
			vals = append(vals, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7))
			args = append(args, row.IPAddress)
			args = append(args, row.CountryCode)
			args = append(args, row.Country)
			args = append(args, row.City)
			args = append(args, row.Latitude)
			args = append(args, row.Longitude)
			args = append(args, row.MysteryValue.(string))
			i++
		}
		qry := fmt.Sprintf("INSERT INTO geolocation (ip_address, country_code, country, city, latitude, longitude, mystery_value) VALUES %s ON CONFLICT (ip_address) DO NOTHING",
			strings.Join(vals, ","))
		_, err := s.db.ExecContext(ctx, qry, args...)
		if err != nil {
			return err
		}
	}
	success = true
	return nil
}

// Get returns a geolocation data matching the give ip address
func (s *Connection) Get(ctx context.Context, ip string) (*model.Location, error) {
	if s.closed {
		return nil, utils.ErrInvalidConn
	}
	qry := fmt.Sprintf("SELECT country_code, country, city, latitude, longitude, mystery_value FROM geolocation WHERE ip_address = '%v'", ip)
	rows, err := s.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	location := model.Location{}
	if rows.Next() {
		err = rows.Scan(&location.CountryCode, &location.Country, &location.City, &location.Latitude, &location.Longitude, &location.MysteryValue)
		if err != nil {
			return nil, errors.New("error scanning detail of ip: " + location.IPAddress + ", err: " + err.Error())
		}
		return &location, nil
	}
	return nil, utils.ErrNotFound
}

// Close closes the database connection
func (s *Connection) Close() error {
	if s.closed {
		return utils.ErrInvalidConn
	}
	return s.db.DB.Close()
}
