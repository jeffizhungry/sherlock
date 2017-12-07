package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jeffizhungry/sherlock/pkg/debug"
)

// TransparentProxy proxies HTTP requests as a MITM service.
type TransparentProxy struct{}

// NewTransparentProxy creates a new TransparentProxy object.
func NewTransparentProxy() *TransparentProxy {
	return &TransparentProxy{}
}

// ServerHTTP handles proxying HTTP requests
func (p *TransparentProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	debug.PPrintln(r)
}

func main() {
	fmt.Println("Sherlock starting")
	defer fmt.Println("Sherlock exiting...")

	server := NewTransparentProxy()
	log.Fatal(http.ListenAndServe(":9999", server))
}
