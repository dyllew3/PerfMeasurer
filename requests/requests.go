package requests

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dyllew3/PerfMeasurer/parser"
)

// Hit list of endpoints specified in Requests file
func HitEndpoints(requestsFile parser.RequestsFile) {
	var host string = requestsFile.Address
	var baseHeaders map[string]string = requestsFile.Headers
	for _, endpoint := range requestsFile.Endpoints {
		var start = time.Now()
		response, err := HitEndpoint(host, baseHeaders, endpoint)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(response)
		var end = time.Now()
		fmt.Println(end.Sub(start).Milliseconds())
	}
}

// Hit specific endpoint
func HitEndpoint(host string, baseHeaders map[string]string, endpoint parser.Endpoint) (*http.Response, error) {
	if !strings.HasSuffix(host, "/") {
		host = host[:len(host)-1]
	}
	url := endpoint.Url
	if !strings.HasPrefix("/", url) {
		url = "/" + url
	}

	var fullPath string = host + url

	client := &http.Client{}
	req, _ := http.NewRequest(endpoint.Method, fullPath, nil)
	// Setting base headers
	for key, val := range baseHeaders {
		req.Header.Set(key, val)
	}

	// Setting headers for specific endpo
	if endpoint.Headers != nil {
		// Setting endpoint specific headers
		for key, val := range *endpoint.Headers {
			req.Header.Set(key, val)
		}
	}
	return client.Do(req)
}
