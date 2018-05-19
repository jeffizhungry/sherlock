package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jeffizhungry/sherlock/microapp"
	"github.com/jeffizhungry/sherlock/pkg/rawhttp"
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

	// Start test http server
	go basicHTTP()

	pairsc := make(chan rawhttp.Pair)

	// Start proxy
	go microapp.Proxy(flagPort, pairsc)

	// Start accumulator
	microapp.Accumulator(pairsc)

	// for p := range sub {
	// 	fmt.Println("Request:  ", p.Request.Method+" "+p.Request.URL.String())
	// 	fmt.Println("Response: ", p.Response.Status)
	// }
}

func basicHTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\n  \"code\": 200,\n  \"description\": \"Hello World\"\n}")
}

func basicHTTP() {
	http.HandleFunc("/", basicHTTPHandler)
	fmt.Println("Starting Basic HTTP Server. Listening on localhost:8080")
	defer fmt.Println("Basic HTTP Server exiting...")
	_ = http.ListenAndServe(":8080", nil)
}
