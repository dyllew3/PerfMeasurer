package generation

import "time"

// Struct representing
type Output struct {
	Title   string
	Host    string
	Date    time.Time
	Results []RequestsResult
}

type RequestsResult struct {
	Endpoint     string
	ResponseTime int64
	StatusCode   int32
}
