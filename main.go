package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Sirupsen/logrus"
)

var (
	flagPort     string
	flagSSL      bool
	flagCertFile string
	flagKeyFile  string
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	flag.StringVar(&flagPort, "port", "9999", "Specify port number to run proxy service")
	flag.BoolVar(&flagSSL, "ssl", false, "Serve as an HTTPS proxy")
	flag.StringVar(&flagCertFile, "cert", "cert.pem", "Specify location of certificate")
	flag.StringVar(&flagKeyFile, "key", "key.pem", "Specify location of private key")
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

	if flagSSL {
		log.Fatal(http.ListenAndServeTLS(":"+flagPort, flagCertFile, flagKeyFile, server))
	}
	log.Fatal(http.ListenAndServe(":"+flagPort, server))
}
