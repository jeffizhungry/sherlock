package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jeffizhungry/sherlock/pkg/debug"
	"github.com/jeffizhungry/sherlock/pkg/rawhttp"
)

var WhitelistDomains = []string{}

// Accumulator is responsible for taking rawhttp.Pairs and storing them
// in the data directory, apply any filters as necessary
func Accumulator(pairs <-chan rawhttp.Pair) {
	log := logrus.WithFields(logrus.Fields{
		"context": "accumulator",
	})
	for p := range pairs {

		// Skip non json requests
		if p.Request.Header.Get("Content-Type") != "application/json" {
			continue
		}

		// Parse HTTP request, skip if it's not empty and invalid

		// Parse HTTP response, skip if it's not empty and invalid

		// Apply whitelist filter
		if len(WhitelistDomains) > 0 {
			host := p.Request.URL.Hostname()
			var found bool
			for _, d := range WhitelistDomains {
				if d == host {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Save
		if err := savePair(p); err != nil {
			log.WithFields(logrus.Fields{
				"method":   p.Request.Method,
				"hostname": p.Request.URL.Hostname(),
				"path":     p.Request.URL.Path,
			}).WithError(err).Error("Failed to save HTTP pair")
		}
	}
}

func savePair(p rawhttp.Pair) error {

	// Make hostname directory
	dirname := "./data/" + p.Request.URL.Hostname()
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return err
	}

	// Save HTTP request
	// filename := dir + "/" +  call.Timestamp.Format(time.RFC3339Nano) + "_random"
	// if err := ioutil.WriteFile(file, data, 0664); err != nil {
	// 	return err
	// }

	// Save HTTP response
	// filename := dir + "/" +  call.Timestamp.Format(time.RFC3339Nano) + "_random"
	// if err := ioutil.WriteFile(file, data, 0664); err != nil {
	// 	return err
	// }
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
	debug.PPrintln("request: ", string(requestPayload))
	if err := json.Unmarshal(requestPayload, &call.RequestBody); err != nil && err != io.EOF {
		debug.PPrintln(req)
		// Ignore non JSON requests
		return nil
	}

	responsePayload, _ := ioutil.ReadAll(resp.Body)
	debug.PPrintln("reponse: ", string(responsePayload))
	if err := json.Unmarshal(responsePayload, &call.ResponseBody); err != nil && err != io.EOF {
		debug.PPrintln(req)
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
