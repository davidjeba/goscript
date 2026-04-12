package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidjeba/goscript/pkg/goscript"
	"github.com/davidjeba/goscript/pkg/components"
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
	html := components.Home(nil)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func helloHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello from GoScript API!")
}

