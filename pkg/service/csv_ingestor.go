package service

import (
	"context"
	"encoding/csv"
	"errors"
	"geolocation/internal"
	"geolocation/internal/model"
	"geolocation/internal/utils"
	"io"
	"strings"
	"time"
)

// CSVIngestor is intended to read valid geolocation data from a csv file
// and ingest it in database
type CSVIngestor struct {
	store internal.Store
	r     io.Reader
}

// NewCSVIngestor creates new instance of CSVIngestor
func NewCSVIngestor(store internal.Store, r io.Reader) *CSVIngestor {
	return &CSVIngestor{store, r}
}

// Ingest reads the csv file
// extracts only valid location data by sanitising it.
// Returns statistics in form of accepted & rejected locations and total time taken, alog with all valid locations
// possibly an error if any step fails
func (c *CSVIngestor) Ingest(ctx context.Context) (*model.Stat, []model.Location, error) {
	// make sure dependencies are intact or fail fast
	if c.store == nil {
		return nil, nil, errors.New("nil store")
	}
	if c.r == nil {
		return nil, nil, errors.New("nil reader")
	}

	now := time.Now()
	s := model.Stat{}
	locations := []model.Location{}
	csvRdr := csv.NewReader(c.r)
	i := 0

	// read & prepare valid locations
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}
		values, err := csvRdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		if i == 0 { // skip header
			header := strings.Join(values, ",")
			if header != csv_header {
				return nil, nil, errors.New("invalid csv header, it must eqauls: " + csv_header)
			}
			i++
			continue
		}
		if row := sanitise(values); row != nil {
			locations = append(locations, *row)
			s.Accepted++
		} else {
			s.Discarded++
		}
	}

	// ingest locations in database
	err := c.store.BulkCreate(ctx, locations)
	if err != nil {
		return nil, nil, errors.New("BulkCreate() failed, err: " + err.Error())
	}

	s.TimeSpent = time.Since(now)
	return &s, locations, nil
}

func sanitise(values []string) *model.Location {
	location := model.Location{}
	for i, value := range values {
		switch column(i) {
		case ip_address:
			if valid, ip := utils.IsIPValid(value); !valid {
				return nil
			} else {
				location.IPAddress = ip
			}
		case country_code:
			if valid, cc := utils.IsStringValid(value); !valid {
				return nil
			} else {
				location.CountryCode = cc
			}
		case country:
			if valid, c := utils.IsStringValid(value); !valid {
				return nil
			} else {
				location.Country = c
			}
		case city:
			if valid, c := utils.IsStringValid(value); !valid {
				return nil
			} else {
				location.City = c
			}
		case latitude:
			if valid, l := utils.IsFloat64Valid(value); !valid {
				return nil
			} else {
				location.Latitude = l
			}
		case longitude:
			if valid, l := utils.IsFloat64Valid(value); !valid {
				return nil
			} else {
				location.Longitude = l
			}
		case mystery_value:
			location.MysteryValue = value
		}
	}
	return &location
}

type column uint

const (
	ip_address column = iota
	country_code
	country
	city
	latitude
	longitude
	mystery_value
)
const csv_header = "ip_address,country_code,country,city,latitude,longitude,mystery_value"
