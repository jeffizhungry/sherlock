package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// TransparentProxy proxies HTTP requests as a MITM service.
type TransparentProxy struct {
	log *logrus.Entry
}

// NewTransparentProxy creates a new TransparentProxy object.
func NewTransparentProxy() *TransparentProxy {
	return &TransparentProxy{
		log: logrus.WithField("context", "proxy"),
	}
}

// ServerHTTP handles proxying HTTP requests
func (p *TransparentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &http.Client{}

	copiedReq, err := copyHTTPRequest(r)
	if err != nil {
		p.log.WithError(err).Error("Failed to copy HTTP request")
		return
	}

	// Make request
	resp, err := c.Do(copiedReq)
	if err != nil {
		p.log.WithError(err).Error("Failed to run HTTP request")
		return
	}

	// Copy and relay response
	copyHTTPHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		p.log.WithError(err).Error("Failed to return HTTP response")
		return
	}
}

func copyHTTPRequest(r *http.Request) (*http.Request, error) {
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if len(bodyData) > 0 {
		body = bytes.NewBuffer(bodyData)
	}
	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		return nil, err
	}
	copyHTTPHeaders(r.Header, req.Header)
	return req, nil
}

func copyHTTPHeaders(src, dst http.Header) {
	for header, values := range src {
		for _, value := range values {
			dst.Add(header, value)
		}
	}
}
