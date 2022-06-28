package model

type Location struct {
	IPAddress    string      `json:"ip_address"`
	CountryCode  string      `json:"country_code"`
	Country      string      `json:"country"`
	City         string      `json:"city"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	MysteryValue interface{} `json:"mystery_value"`
}
