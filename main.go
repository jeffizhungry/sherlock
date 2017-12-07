package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jeffizhungry/sherlock/pkg/debug"
)

// TransparentProxy proxies HTTP requests as a MITM service.
type TransparentProxy struct{}

// NewTransparentProxy creates a new TransparentProxy object.
func NewTransparentProxy() *TransparentProxy {
	return &TransparentProxy{}
}

// ServerHTTP handles proxying HTTP requests
func (p *TransparentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &http.Client{}

	debug.PPrintln("request:", r)

	// Copy request
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		debug.Println("tproxy: read body error: ", err.Error())
		return
	}
	var body io.Reader
	if len(bodyData) > 0 {
		body = bytes.NewBuffer(bodyData)
	}
	req, err := http.NewRequest(r.Method, r.URL.String(), body)
	if err != nil {
		debug.Println("tproxy: make new request: ", err.Error())
		return
	}
	copyHeader(r.Header, req.Header)

	// Make request
	resp, err := c.Do(req)
	if err != nil {
		debug.Println("tproxy: http request error: ", err.Error())
		return
	}

	debug.PPrintln("response:", r)

	// Copy and relay response
	copyHeader(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		debug.Println("tproxy: http response copy error: ", err.Error())
		return
	}
}

func copyHeader(src, dst http.Header) {
	for header, values := range src {
		for _, value := range values {
			dst.Add(header, value)
		}
	}
}

func main() {
	fmt.Println("Sherlock starting")
	defer fmt.Println("Sherlock exiting...")

	server := NewTransparentProxy()
	log.Fatal(http.ListenAndServe(":9999", server))
}
