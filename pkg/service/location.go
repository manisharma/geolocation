package service

import (
	"encoding/json"
	"errors"
	"geolocation/internal"
	"geolocation/internal/model"
	"geolocation/internal/utils"
	"net/http"

	"github.com/sirupsen/logrus"
)

type LocationService struct {
	store internal.Store
}

// New LocationService creates new instance of LocationService
func NewLocationService(store internal.Store) *LocationService {
	return &LocationService{store}
}

// Resolve handles ip resolution request
func (l *LocationService) Resolve(w http.ResponseWriter, r *http.Request) {
	// onyl GET & OPTION method are allowed
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// read ip address from query params
	ip := r.URL.Query().Get("ip")

	// validate ip
	if ok, _ := utils.IsIPValid(ip); !ok {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		utils.GetLogger().WithFields(logrus.Fields{"ip": ip}).Error(utils.ErrBadRequest.Error())
		w.Write([]byte(utils.ErrBadRequest.Error()))
		return
	}

	// initialise a common logger
	logger := utils.GetLogger().WithFields(logrus.Fields{"ip": ip})

	// resolve the ip or fail fast
	row, err := l.store.Get(r.Context(), ip)
	if err != nil {
		logger.WithFields(logrus.Fields{"err": err}).Error("resolving ip address failed")
		w.Header().Set("Content-Type", "text/html")
		if errors.Is(err, utils.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(utils.ErrNotFound.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(utils.ErrInternalServerError.Error()))
		return
	}

	// prepare response
	logger.WithFields(logrus.Fields{"row": row}).Debug("ip address resolved")
	bytes, err := json.MarshalIndent(model.Location{
		IPAddress:    ip,
		CountryCode:  row.CountryCode,
		Country:      row.Country,
		City:         row.City,
		Latitude:     row.Latitude,
		Longitude:    row.Longitude,
		MysteryValue: row.MysteryValue,
	}, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithFields(logrus.Fields{"err": err}).Error("MarshalIndent() failed")
		w.Write([]byte(utils.ErrInternalServerError.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
