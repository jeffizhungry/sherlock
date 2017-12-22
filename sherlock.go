package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"polymail-api/lib/utils"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jeffizhungry/sherlock/pkg/debug"
)

// Sherlock inspects and deciphers HTTP payloads
type Sherlock struct {
	log      *logrus.Entry
	payloads <-chan HTTPPayload
}

// NewSherlock creates a new Sherlock instnace.
func NewSherlock(payloads <-chan HTTPPayload) *Sherlock {
	return &Sherlock{
		log:      logrus.WithField("context", "proxy"),
		payloads: payloads,
	}
}

// Run reads and processes HTTP paylaods
func (s *Sherlock) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-s.payloads:
			if err := s.inspect(p.Request, p.Response); err != nil {
				s.log.WithError(err).Error("Inspection Error")
			}
		}
	}
}

func (s *Sherlock) inspect(req *http.Request, resp *http.Response) error {
	s.log.Debugf("Inspecting request: %v %v - %v", req.Method, req.RequestURI, resp.StatusCode)

	// Ignore error cases for now
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// debug.PPrintln("BAD STATUS CODE:", resp.StatusCode)
		// data, _ := ioutil.ReadAll(resp.Body)
		// debug.PPrintln("BAD PAYLOAD:", string(data))
		return nil
	}

	// Reject non-JSON payloads
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		// debug.PPrintln("BAD CONTENT TYPE:", contentType)
		return nil
	}

	// Save
	if err := savePayload(req, resp); err != nil {
		return err
	}
	return nil
}

func savePayload(req *http.Request, resp *http.Response) error {

	// Format Data, Ignore non-JSON payloads
	call := APIRouteCall{
		Timestamp:  time.Now(),
		URL:        req.URL.String(),
		StatusCode: resp.StatusCode,
	}

	// Copy payloads
	requestPayload, _ := ioutil.ReadAll(req.Body)
	utils.PPrintln("request: ", string(requestPayload))
	if err := json.Unmarshal(requestPayload, &call.RequestBody); err != nil && err != io.EOF {
		utils.PPrintln(req)
		// Ignore non JSON requests
		return nil
	}

	responsePayload, _ := ioutil.ReadAll(resp.Body)
	utils.PPrintln("reponse: ", string(responsePayload))
	if err := json.Unmarshal(responsePayload, &call.ResponseBody); err != nil && err != io.EOF {
		utils.PPrintln(req)
		// Ignore non JSON responses
		return nil
	}

	// Copy headers
	for h := range req.Header {
		call.RequestHeaders = append(call.RequestHeaders, h)
	}
	for h := range resp.Header {
		call.ResponseHeaders = append(call.ResponseHeaders, h)
	}

	// Make Directory
	dir := "data/" + req.URL.Hostname() + "/" + req.URL.Path
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Save
	debug.PPrintln("saving: ", call)
	// data, err := json.MarshalIndent(call, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// if err := ioutil.WriteFile(dir + "/payload_" + call.Timestamp.Format(time.RFC3339Nano), data, 0664); err != nil {
	// 	return err
	// }
	return nil
}

// func (s *Sherlock) print(req *http.Request, resp *http.Response) error {
//
// 	// Print Request
// 	debug.PPrintln("request.URL:", req.URL.String())
// 	requestBody, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		logrus.WithError(err).Error("Failed to read request body")
// 	} else {
// 		debug.PPrintln("request.Body:", string(requestBody))
// 	}
//
// 	// Print Response
// 	debug.PPrintln("response.Status:", resp.Status)
// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		logrus.WithError(err).Error("Failed to read resposne body")
// 	} else {
// 		debug.PPrintln("response.Body:", string(responseBody))
// 	}
// 	return nil
// }
//
