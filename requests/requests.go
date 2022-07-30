package requests

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dyllew3/PerfMeasurer/generation"
	"github.com/dyllew3/PerfMeasurer/parser"
)

// Hit list of endpoints specified in Requests file
func HitEndpoints(requestsFile parser.RequestsFile) generation.Output {
	log.Printf("Performing requests on host %s\n", requestsFile.Address)

	var host string = requestsFile.Address
	var outputResult generation.Output = generation.Output{}
	outputResult.Date = time.Now()
	outputResult.Title = requestsFile.Name
	outputResult.Host = requestsFile.Address

	var baseHeaders map[string]string = requestsFile.Headers

	// Go through each endpoint and fire off a request
	for _, endpoint := range requestsFile.Endpoints {
		log.Printf("Hitting endpoint %s\n", endpoint.Url)
		var start = time.Now()
		response, err := HitEndpoint(host, baseHeaders, endpoint)
		var end = time.Now()

		// Format endpoint result
		if err != nil {
			log.Printf("Enocuntered error at endpoint %s", endpoint.Url)
			outputResult.Results = append(outputResult.Results, generation.RequestsResult{
				Endpoint:       endpoint.Url,
				ResponseTimeMs: end.Sub(start).Milliseconds(),
				Method:         endpoint.Method,
			})
		} else {
			outputResult.Results = append(outputResult.Results, generation.RequestsResult{
				StatusCode:     int32(response.StatusCode),
				Endpoint:       endpoint.Url,
				ResponseTimeMs: end.Sub(start).Milliseconds(),
				Method:         endpoint.Method,
			})
		}
	}
	return outputResult
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
	req, _ := http.NewRequest(endpoint.Method, fullPath, strings.NewReader(endpoint.Body))
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
