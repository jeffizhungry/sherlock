package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"polymail-api/lib/utils"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// HTTPPayload encapsulates a HTTP request and response pair
type HTTPPayload struct {
	Request  *http.Request
	Response *http.Response
}

// TransparentProxy proxies HTTP requests as a MITM service.
type TransparentProxy struct {
	log     *logrus.Entry
	payload chan<- HTTPPayload
	client  *http.Client
}

// NewTransparentProxy creates a new TransparentProxy object.
func NewTransparentProxy(p chan<- HTTPPayload) *TransparentProxy {
	return &TransparentProxy{
		log:     logrus.WithField("context", "proxy"),
		payload: p,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// ServerHTTP handles proxying HTTP requests
func (p *TransparentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (r.Host == "localhost:"+flagPort || r.Host == "127.0.0.1:"+flagPort) && r.URL.Path == "/" {
		msg := `Proxy is alive and running!`
		p.log.Info(msg)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(msg))
		return
	}

	// Copy Request
	copiedReq, err := copyHTTPRequest(r)
	if err != nil {
		p.log.WithError(err).Error("Failed to copy HTTP request")
		return
	}
	// p.log.Debugf("Proxying request: %v", copiedReq.RequestURI)

	// Make request
	resp, err := p.client.Do(copiedReq)
	if err != nil {
		utils.PPrintln(r.URL, r.Header)
		p.log.WithError(err).Error("Failed to run HTTP request")
		return
	}

	// Copy response
	if err := copyHTTPResponse(resp, w); err != nil {
		p.log.WithError(err).Error("Failed to return HTTP response")
		return
	}

	// Forward request response pairs if channel is not nil
	select {
	case p.payload <- HTTPPayload{Request: r, Response: resp}:
	default:
	}
}

func copyHTTPRequest(r *http.Request) (*http.Request, error) {
	var copiedBody io.Reader
	if r.ContentLength > 0 && r.Body != nil {
		var buf bytes.Buffer
		copiedBody = io.TeeReader(r.Body, &buf)
		r.Body = ioutil.NopCloser(&buf)
	}
	url := r.URL.String()
	if r.URL.Scheme == "" && strings.HasPrefix(url, "//") {
		if origin := r.Header.Get("Origin"); origin != "" {
			originSchemeURL := strings.Split(origin, "://")
			url = originSchemeURL[0] + ":" + url
		} else {
			// Default to HTTPS protocol. Should be good for most cases
			url = "https:" + url
		}
	}
	copied, err := http.NewRequest(r.Method, url, copiedBody)
	if err != nil {
		return nil, err
	}
	copyHTTPHeaders(r.Header, copied.Header)
	return copied, nil
}

func copyHTTPResponse(resp *http.Response, w http.ResponseWriter) error {
	copyHTTPHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)

	if resp.ContentLength > 0 {
		var buf bytes.Buffer
		copiedBody := io.TeeReader(resp.Body, &buf)
		resp.Body = ioutil.NopCloser(&buf)
		_, err := io.Copy(w, copiedBody)
		if err != nil {
			return err
		}
	} else {
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyHTTPHeaders(src, dst http.Header) {
	for header, values := range src {
		for _, value := range values {
			dst.Add(header, value)
		}
	}
}

// bodyData, err := ioutil.ReadAll(r.Body)
// if err != nil {
// 	return nil, err
// }
// var body io.Reader
// if len(bodyData) > 0 {
// 	body = bytes.NewBuffer(bodyData)
// }
