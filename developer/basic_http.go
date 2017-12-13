package main

import (
	"fmt"
	"net/http"

	"github.com/jeffizhungry/sherlock/pkg/debug"
)

func handler(w http.ResponseWriter, r *http.Request) {
	debug.PPrintln("request:", r)
	fmt.Fprintf(w, "{\n  \"code\": 200,\n  \"description\": \"Hello World\"\n}")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting Basic HTTP Server. Listening on localhost:8080")
	defer fmt.Println("Basic HTTP Server exiting...")
	_ = http.ListenAndServe(":8080", nil)
}
