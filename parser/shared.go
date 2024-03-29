package parser

// Struct used to represent yaml file
// Has the necessary
type RequestsFile struct {
	Name                 string            `yaml:"name" json:"name"`
	ReportOutputFilename *string           `yaml:"output_name" json:"output_name"`
	Address              string            `yaml:"host" json:"host"`
	Headers              map[string]string `yaml:"headers" json:"headers"`
	Endpoints            []Endpoint        `yaml:"endpoints" json:"endpoints"`
	NumRequests          *int64            `yaml:"num_requests" json:"num_requests"`
	FollowRobots         *bool             `yaml:"follow_robots" json:"follow_robots"`
}

type Endpoint struct {
	Url         string             `yaml:"url" json:"url"`
	Method      string             `yaml:"method" json:"method"`
	Body        string             `yaml:"body" json:"body"`
	BodyParams  *map[string]string `yaml:"body_params" json:"body_params"`
	UrlParams   *map[string]string `yaml:"url_params" json:"url_params"`
	Timeout     *int32             `yaml:"timeout" json:"timeout"`
	NumRequests *int64             `yaml:"num_requests" json:"num_requests"`
	Headers     *map[string]string `yaml:"headers" json:"headers"`
}
