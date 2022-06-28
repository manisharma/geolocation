package model

import "time"

type Stat struct {
	TimeSpent time.Duration `json:"timeSpent"`
	Accepted  int           `json:"accepted"`
	Discarded int           `json:"discarded"`
}
