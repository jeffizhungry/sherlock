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
	_ = http.ListenAndServe(":8080", nil)
}
