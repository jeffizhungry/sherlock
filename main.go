package main

import (
	"fmt"
	"log"
	"net/http"
)

// Proxy proxies HTTP requests
type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("Sherlock starting")
	defer fmt.Println("Sherlock exiting...")
	server := &Proxy{}
	log.Fatal(http.ListenAndServe(":9999", server))
}
