package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var flagPort string

func init() {
	flag.StringVar(&flagPort, "port", "9999", "Specify port number to run proxy service")
	flag.Parse()
}

func main() {
	fmt.Println("Sherlock starting. Listening on localhost:" + flagPort)
	defer fmt.Println("Sherlock exiting...")

	// Create channels
	payloads := make(chan HTTPPayload)

	// Start Consumer
	sher := NewSherlock(payloads)
	go sher.Run(context.TODO())

	// Start Proxy
	server := NewTransparentProxy(payloads)
	log.Fatal(http.ListenAndServe(":"+flagPort, server))
}
