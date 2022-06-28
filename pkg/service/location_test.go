package service

import (
	"context"
	"geolocation/internal"
	"geolocation/internal/model"
	"geolocation/internal/store/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLocationService(t *testing.T) {
	conn, err := mock.New()
	assert.Nil(t, err, "NewMockConnection() failed, expected no error, got %v", err)
	type args struct {
		store internal.Store
	}
	tests := []struct {
		name string
		args args
		want *LocationService
	}{
		{
			name: "should create new instance",
			args: args{
				store: conn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLocationService(tt.args.store)
			assert.NotNil(t, got, "NewLocationService() failed, expected not nil, got nil")
		})
	}
}

func TestLocationService_Resolve(t *testing.T) {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 0*time.Second)
	defer cancel()
	conn, err := mock.New()
	assert.Nil(t, err, "NewMockConnection() failed, expected no error, got %v", err)
	type fields struct {
		store internal.Store
	}
	tests := []struct {
		name   string
		fields fields
		ip     string
		status int
		seed   *model.Location
		method string
		ctx    context.Context
		ctxErr bool
	}{
		{
			name:   "valid ip, should be resolved with status of 200 OK",
			fields: fields{conn},
			ip:     "192.168.0.0",
			status: http.StatusOK,
			seed: &model.Location{
				IPAddress:    "192.168.0.0",
				CountryCode:  "IN",
				Country:      "India",
				City:         "Bengaluru",
				Latitude:     19.3445466755,
				Longitude:    45.9454878475,
				MysteryValue: "anything",
			},
			method: http.MethodGet,
			ctx:    ctx,
			ctxErr: false,
		},
		{
			name:   "missing ip, should be resolved with status of 400 Bad Request",
			fields: fields{conn},
			ip:     "",
			status: http.StatusBadRequest,
			seed:   nil,
			method: http.MethodGet,
			ctx:    ctx,
			ctxErr: false,
		},
		{
			name:   "invalid ip, should be resolved with status of 400 Bad Request",
			fields: fields{conn},
			ip:     "192.168.abc.cba",
			status: http.StatusBadRequest,
			seed:   nil,
			method: http.MethodGet,
			ctx:    ctx,
			ctxErr: false,
		},
		{
			name:   "an ip, which has not been ingested should be resolved with status of 404 Not Found",
			fields: fields{conn},
			ip:     "125.159.20.54",
			status: http.StatusNotFound,
			seed:   nil,
			method: http.MethodGet,
			ctx:    ctx,
			ctxErr: false,
		},
		{
			name:   "any other error (catch all) should be resolved with status of 500 Internal Server Error",
			fields: fields{conn},
			ip:     "125.159.20.54",
			status: http.StatusInternalServerError,
			seed: &model.Location{
				IPAddress:    "192.168.0.0",
				CountryCode:  "IN",
				Country:      "India",
				City:         "Bengaluru",
				Latitude:     19.3445466755,
				Longitude:    45.9454878475,
				MysteryValue: "1234567775",
			},
			method: http.MethodGet,
			ctx:    ctxWithTimeout,
			ctxErr: true,
		},
		{
			name:   "except GET or OPTION requests no other requests are allowed, should be resolved with status of 405 Method Not Allowed",
			fields: fields{conn},
			ip:     "125.159.20.54",
			status: http.StatusMethodNotAllowed,
			seed:   nil,
			method: http.MethodPost,
			ctx:    ctx,
			ctxErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(tt.method, "/resolve?"+url.Values{"ip": {tt.ip}}.Encode(), nil)
			r = r.WithContext(tt.ctx)
			assert.Nil(t, err, "NewRequest() failed, expected no error, got %v", err)
			l := &LocationService{tt.fields.store}
			if tt.seed != nil {
				err = l.store.BulkCreate(r.Context(), []model.Location{*tt.seed})
				if tt.ctxErr {
					assert.Equal(t, context.DeadlineExceeded, err, "BulkCreate() expecting error %v, got %v", context.DeadlineExceeded, err)
				} else {
					assert.Nil(t, err, "BulkCreate() failed, expected no error, got %v", err)
				}
			}
			l.Resolve(w, r)
			assert.Equal(t, tt.status, w.Code, "expectd status %d, got %d", tt.status, w.Code)
		})
	}
}
