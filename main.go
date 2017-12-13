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
	flagSSLPort  string
	flagSSL      bool
	flagCertFile string
	flagKeyFile  string
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	flag.StringVar(&flagPort, "port", "9090", "Specify port number to run HTTP proxy service")
	flag.StringVar(&flagSSLPort, "sslport", "9091", "Specify port number to run HTTPS proxy service")
	flag.BoolVar(&flagSSL, "ssl", false, "Serve as an HTTPS proxy")
	flag.StringVar(&flagCertFile, "cert", "cert.pem", "Specify location of certificate")
	flag.StringVar(&flagKeyFile, "key", "key.pem", "Specify location of private key")
	flag.Parse()
}

func main() {
	defer fmt.Println("Sherlock exiting...")

	// Create channels
	payloads := make(chan HTTPPayload)

	// Start Consumer
	sher := NewSherlock(payloads)
	go sher.Run(context.TODO())

	// Start HTTP and HTTPS Proxy
	server := NewTransparentProxy(payloads)
	go func() {
		fmt.Println("Sherlock HTTPS Proxy. Listening on localhost:" + flagSSLPort)
		log.Fatal(http.ListenAndServeTLS("localhost:"+flagSSLPort, flagCertFile, flagKeyFile, server))
	}()
	fmt.Println("Sherlock HTTS Proxy. Listening on localhost:" + flagPort)
	log.Fatal(http.ListenAndServe("localhost:"+flagPort, server))
}
