package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
	router := goscript.NewRouter()

	// Register routes
	router.GET("/", homeHandler)
	router.GET("/api/hello", helloHandler)

	// Start the server
	port := 8080
	fmt.Printf("Server starting on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func homeHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Welcome to GoScript!")
}

func helloHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello from GoScript API!")
}

