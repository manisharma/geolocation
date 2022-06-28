package utils

import (
	"errors"
	"geolocation/internal/model"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger *log.Logger

const (
	CreateTableScript = `CREATE TABLE IF NOT EXISTS geolocation (
		ip_address 		varchar not null PRIMARY KEY,
		country_code 	varchar not null,
		country 		varchar not null,
		city 			varchar not null,
		latitude 		varchar not null,
		longitude 		varchar not null,
		mystery_value 	varchar not null)`
)

var (
	ErrInvalidConn         = errors.New("invalid db connection")
	ErrBadRequest          = errors.New("you must provide a valid ip address")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("there was an error processing your request, please re-try after sometime")
)

// GetLogger returns singleton logger
func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New()
		logger.SetLevel(log.DebugLevel)
	}
	return logger
}

// GetDBCfg prepares database configuration by reading ENV vars
func GetDBCfg() (*model.DB, error) {
	dbCfg := model.DB{}
	err := viper.Unmarshal(&dbCfg)
	if err != nil {
		return nil, err
	}
	return &dbCfg, nil
}

func IsIPValid(ip string) (bool, string) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false, ip
	}
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return false, ip
		}
	}
	return true, ip
}

func IsStringValid(s string) (bool, string) {
	s = strings.TrimSpace(s)
	return len(s) > 0, s
}

func IsFloat64Valid(s string) (bool, float64) {
	v, err := strconv.ParseFloat(s, 64)
	return err == nil, v
}
