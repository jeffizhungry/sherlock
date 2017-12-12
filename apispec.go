package main

import "time"

// APIRouteCall captures basic information from an API call
type APIRouteCall struct {
	Timestamp       time.Time
	URL             string
	StatusCode      int
	RequestBody     map[string]interface{}
	RequestHeaders  []string
	ResponseBody    map[string]interface{}
	ResponseHeaders []string
}
