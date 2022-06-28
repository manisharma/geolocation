package service

import (
	"context"
	"geolocation/internal"
	"geolocation/internal/model"
	"geolocation/internal/store/mock"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVReader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *CSVIngestor
	}{
		{
			name: "should create new instance",
			args: args{
				r: strings.NewReader(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCSVIngestor(nil, tt.args.r)
			assert.NotNil(t, got, "NewCSVReader() failed, expected not nil, got nil")
		})
	}
}

func TestCSVReader_Read(t *testing.T) {
	type fields struct {
		r io.Reader
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Stat
		want1   []model.Location
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CSVIngestor{
				r: tt.fields.r,
			}
			got, got1, err := c.Ingest(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVReader.Read() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CSVReader.Read() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_sanitise(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		result     *model.Location
		shouldPass bool
	}{
		{
			name: "valid location should pass",
			line: "200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346",
			result: &model.Location{
				IPAddress:    "200.106.141.15",
				CountryCode:  "SI",
				Country:      "Nepal",
				City:         "DuBuquemouth",
				Latitude:     -84.87503094689836,
				Longitude:    7.206435933364332,
				MysteryValue: 7823011346,
			},
			shouldPass: true,
		},
		{
			name:       "invalid valid location (ip address) should not pass",
			line:       "200.106.141.tyr,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346",
			result:     nil,
			shouldPass: false,
		},
		{
			name:       "invalid valid location (missing ip address) should not pass",
			line:       ",SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346",
			result:     nil,
			shouldPass: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := strings.Split(tt.line, ",")
			got := sanitise(args)
			if tt.shouldPass {
				assert.NotNil(t, got, "wanted not nil location %#v, got nil location", tt.result)
				assert.Equal(t, tt.result.City, got.City, "city must match, wanted %v, got %v", tt.result.City, got.City)
				assert.Equal(t, tt.result.Country, got.Country, "country must match, wanted %v, got %v", tt.result.Country, got.Country)
				assert.Equal(t, tt.result.CountryCode, got.CountryCode, "country code must match, wanted %v, got %v", tt.result.CountryCode, got.CountryCode)
				assert.Equal(t, tt.result.IPAddress, got.IPAddress, "ip address must match, wanted %v, got %v", tt.result.IPAddress, got.IPAddress)
				assert.Equal(t, tt.result.Latitude, got.Latitude, "latitude must match, wanted %v, got %v", tt.result.Latitude, got.Latitude)
				assert.Equal(t, tt.result.Longitude, got.Longitude, "longitude must match, wanted %v, got %v", tt.result.Longitude, got.Longitude)
			} else {
				assert.Nil(t, got, "wanted nil location, got not nil location", got)
			}
		})
	}
}

func TestCSVIngestor_Ingest(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 0*time.Second)
	defer cancel()
	conn, err := mock.New()
	assert.Nil(t, err, "NewMockConnection() failed, expected no error, got %v", err)
	type dependencies struct {
		store internal.Store
		r     io.Reader
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		dependencies dependencies
		args         args
		stat         *model.Stat
		locations    []model.Location
		wantErr      bool
	}{
		{
			name: "nil store should fail",
			dependencies: dependencies{
				store: nil,
				r:     nil,
			},
			args: args{
				ctx: ctx,
			},
			stat:      nil,
			locations: nil,
			wantErr:   true,
		},
		{
			name: "nil reader should fail",
			dependencies: dependencies{
				store: nil,
				r:     nil,
			},
			args: args{
				ctx: ctx,
			},
			stat:      nil,
			locations: nil,
			wantErr:   true,
		},
		{
			name: "missing csv header should fail",
			dependencies: dependencies{
				store: conn,
				r: strings.NewReader("200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
					"160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n" +
					"70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162\n"),
			},
			args: args{
				ctx: ctx,
			},
			stat:      nil,
			locations: nil,
			wantErr:   true,
		},
		{
			name: "context gets cancelled should fail",
			dependencies: dependencies{
				store: conn,
				r: strings.NewReader("200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
					"160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n" +
					"70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162\n"),
			},
			args: args{
				ctx: ctxWithTimeout,
			},
			stat:      nil,
			locations: nil,
			wantErr:   true,
		},
		{
			name: "valid csv data should pass",
			dependencies: dependencies{
				store: conn,
				r: strings.NewReader("ip_address,country_code,country,city,latitude,longitude,mystery_value\n" +
					"200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n" +
					"160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n" +
					"70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162\n"),
			},
			args: args{
				ctx: ctx,
			},
			stat: &model.Stat{
				Accepted:  3,
				Discarded: 0,
			},
			locations: []model.Location{
				{
					IPAddress:    "200.106.141.15",
					CountryCode:  "SI",
					Country:      "Nepal",
					City:         "DuBuquemouth",
					Latitude:     -84.87503094689836,
					Longitude:    7.206435933364332,
					MysteryValue: 7823011346,
				},
				{
					IPAddress:    "160.103.7.140",
					CountryCode:  "CZ",
					Country:      "Nicaragua",
					City:         "New Neva",
					Latitude:     -68.31023296602508,
					Longitude:    -37.62435199624531,
					MysteryValue: 7301823115,
				},
				{
					IPAddress:    "70.95.73.73",
					CountryCode:  "TL",
					Country:      "Saudi Arabia",
					City:         "Gradymouth",
					Latitude:     -49.16675918861615,
					Longitude:    -86.05920084416894,
					MysteryValue: 2559997162,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CSVIngestor{
				store: tt.dependencies.store,
				r:     tt.dependencies.r,
			}
			stat, locations, err := c.Ingest(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVIngestor.Ingest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.NotNil(t, err, "expected not nill err, got nil")
				assert.Nil(t, stat, "expected nill stat, got %#v", stat)
				assert.Nil(t, locations, "expected nill locations, got %#v", locations)
			} else {
				assert.Equal(t, tt.stat.Accepted, stat.Accepted, "expected accepted location count %d does not match actual accepted location count %d", tt.stat.Accepted, stat.Accepted)
				assert.Equal(t, tt.stat.Discarded, stat.Discarded, "expected discarded location count %d does not match actual discarded location count %d", tt.stat.Discarded, stat.Discarded)
				assert.Equal(t, len(tt.locations), len(locations), "expected location count must match actual location count")
			}
		})
	}
}
