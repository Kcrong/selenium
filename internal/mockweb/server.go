//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	defaultPort = "18080"
)

func main() {
	port := os.Getenv("SELENIUM_TEST_SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir("./htmls")))
	fmt.Printf("Test server running on http://localhost:%s\n", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
