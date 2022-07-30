package generation

import "time"

// Struct representing the final output
type Output struct {
	Title   string           `yaml:"title" json:"title"`
	Host    string           `yaml:"host" json:"host"`
	Date    time.Time        `yaml:"date" json:"date"`
	Results []RequestsResult `yaml:"results" json:"results"`
}

type RequestsResult struct {
	Endpoint       string `yaml:"endpoint" json:"endpoint"`
	ResponseTimeMs int64  `yaml:"response_time_ms" json:"response_time_ms"`
	StatusCode     int32  `yaml:"status_code" json:"status_code"`
	Method         string `yaml:"method" json:"method"`
	Error          *error `yaml:"error" json:"error"`
}
