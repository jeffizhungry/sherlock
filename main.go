package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jeffizhungry/sherlock/pkg/debug"
)

var flagPort string

func init() {
	flag.StringVar(&flagPort, "port", "9999", "Specify port number to run proxy service")
	flag.Parse()
}

func consumer(c <-chan HTTPPayload) {
	for payload := range c {

		// Print Request
		debug.PPrintln("request.URL:", payload.Request.URL.String())
		requestBody, err := ioutil.ReadAll(payload.Request.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read request body")
		} else {
			debug.PPrintln("request.Body:", string(requestBody))
		}

		// Print Response
		debug.PPrintln("response.Status:", payload.Response.Status)
		responseBody, err := ioutil.ReadAll(payload.Response.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read resposne body")
		} else {
			debug.PPrintln("response.Body:", string(responseBody))
		}
	}
}
func main() {
	fmt.Println("Sherlock starting. Listening on localhost:" + flagPort)
	defer fmt.Println("Sherlock exiting...")

	// Create channels
	payload := make(chan HTTPPayload)

	// Start Consumer
	go consumer(payload)

	// Start Proxy
	server := NewTransparentProxy(payload)
	log.Fatal(http.ListenAndServe(":"+flagPort, server))
}
