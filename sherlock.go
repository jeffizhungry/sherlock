package main

import (
	"context"
	"io/ioutil"
	"net/http"

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
	// Print Request
	debug.PPrintln("request.URL:", req.URL.String())
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read request body")
	} else {
		debug.PPrintln("request.Body:", string(requestBody))
	}

	// Print Response
	debug.PPrintln("response.Status:", resp.Status)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read resposne body")
	} else {
		debug.PPrintln("response.Body:", string(responseBody))
	}
	return nil
}
