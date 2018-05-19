package microapp

import (
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jeffizhungry/sherlock/pkg/random"
	"github.com/jeffizhungry/sherlock/pkg/rawhttp"
)

// Accumulator is responsible for taking rawhttp.Pairs and storing them
// in the data directory, apply any filters as necessary. Accumulator
// can be run concurrently as part of a worker pool.
func Accumulator(pairs <-chan rawhttp.Pair) {
	log := logrus.WithFields(logrus.Fields{
		"context": "accumulator",
	})
	for p := range pairs {
		r := p.Request.Method + " " + p.Request.URL.String()
		fmt.Println("Received Raw HTTP Pair: ", r)

		// // Skip non json requests
		// if p.Request.Header.Get("Content-Type") != "application/json" {
		// 	log.WithFields(logrus.Fields{
		// 		"route": r,
		// 	}).Debug("Content-Type is not application/json. Dropping request.")
		// 	continue
		// }

		// // Parse HTTP request, skip if it's not empty and invalid

		// // Parse HTTP response, skip if it's not empty and invalid

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
	dir := filepath.Join(".", "data", p.Request.URL.Hostname())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	prefix := time.Now().Format(time.RFC3339Nano) + random.SafeAlphanumeric(6)

	// Save HTTP request
	reqFile := filepath.Join(dir, prefix+"_request.txt")
	reqData, err := httputil.DumpRequest(p.Request, true)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(reqFile, reqData, 0664); err != nil {
		return err
	}

	// Save HTTP response
	respFile := filepath.Join(dir, prefix+"_response.txt")
	respData, err := httputil.DumpResponse(p.Response, true)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(respFile, respData, 0664); err != nil {
		return err
	}
	return nil
}
