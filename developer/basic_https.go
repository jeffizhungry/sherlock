package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jeffizhungry/sherlock/pkg/debug"
)

func handler(w http.ResponseWriter, r *http.Request) {
	debug.PPrintln("request:", r)
	fmt.Fprintf(w, "{\n  \"code\": 200,\n  \"description\": \"Hello World\"\n}")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting Basic HTTPS Server. Listening on localhost:8080")
	defer fmt.Println("Basic HTTP Server exiting...")
	log.Fatal(http.ListenAndServeTLS("localhost:8080", "cert.pem", "key.pem", nil))
}
