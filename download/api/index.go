package handler

import (
	"net/http"

	"github.com/davidjeba/goscript/pkg/goscript"
	"github.com/davidjeba/goscript/pkg/components"
)

// Handler - Vercel serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	router := goscript.NewRouter()

	// Register routes
	router.GET("/", homeHandler)
	router.GET("/api/hello", helloHandler)

	// Handle the request
	router.ServeHTTP(w, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	html := components.Home(nil)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func helloHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello from GoScript API!"}`))
}

