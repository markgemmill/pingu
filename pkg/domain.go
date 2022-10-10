package pkg

import "time"

// StoreRecord contains ongoing data on a sequential series of the same state.
type StoreRecord struct {
	Start    time.Time `json:"start"`
	Last     time.Time `json:"last"`
	Interval float64   `json:"interval"`
	Count    int64     `json:"count"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}
