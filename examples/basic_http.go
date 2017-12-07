package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\n  \"code\": 200,\n  \"description\": \"Hello World\"\n}")
}

func main() {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":8080", nil)
}
