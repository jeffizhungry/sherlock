package rawhttp

import "net/http"

// Pair represents an HTTP request response pair
type Pair struct {
	Request  *http.Request
	Response *http.Response
}
