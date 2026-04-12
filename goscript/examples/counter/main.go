// Package main is a minimal GoScript 2.0 application — a reactive counter
// built with Go server-side rendering and goscript reactive attributes.
//
// This example demonstrates:
//   - GoScript App Router for page routing
//   - Server-side rendered HTML with zero JavaScript
//   - Reactive attributes (gs-trigger, gs-get, gs-target, gs-swap) for interactivity
//   - API endpoint returning HTML fragments for partial page updates
//   - Metadata/SEO via the goscript builder API
//
// Run:
//
//	go run main.go
//	# or
//	PORT=3000 go run main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	router := goscript.NewAppRouter("/")

	// Home page — renders a reactive counter with server-side HTML
	router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		metadata := goscript.NewMetadata().
			SetTitle("GoScript Counter").
			SetDescription("A reactive counter built with GoScript 2.0 — zero JavaScript").
			SetThemeColor("#10b981").
			Build()

		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  %s
  <script src="/__goscript/runtime.js"></script>
  <style>
    body {
      font-family: 'Inter', system-ui, -apple-system, sans-serif;
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
      margin: 0;
      background: #f9fafb;
    }
    .counter {
      text-align: center;
      padding: 3rem;
      background: white;
      border-radius: 16px;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1), 0 10px 40px rgba(0, 0, 0, 0.05);
      border: 1px solid #f3f4f6;
    }
    .counter h1 {
      color: #111827;
      font-size: 1.5rem;
      margin-bottom: 0.25rem;
    }
    .counter .subtitle {
      color: #6b7280;
      font-size: 0.875rem;
    }
    #count {
      font-size: 5rem;
      font-weight: 800;
      color: #10b981;
      margin: 1.5rem 0;
      line-height: 1;
      transition: transform 0.15s ease;
    }
    .buttons {
      display: flex;
      gap: 0.5rem;
      justify-content: center;
      margin-top: 0.5rem;
    }
    button {
      padding: 0.75rem 1.5rem;
      border: none;
      border-radius: 10px;
      font-size: 1.125rem;
      cursor: pointer;
      transition: all 0.15s ease;
      font-weight: 600;
      min-width: 60px;
    }
    button:active {
      transform: scale(0.95);
    }
    .dec {
      background: #fef2f2;
      color: #dc2626;
      border: 1px solid #fecaca;
    }
    .dec:hover {
      background: #fee2e2;
    }
    .inc {
      background: #ecfdf5;
      color: #059669;
      border: 1px solid #a7f3d0;
    }
    .inc:hover {
      background: #d1fae5;
    }
    .reset {
      background: #f3f4f6;
      color: #374151;
      border: 1px solid #e5e7eb;
    }
    .reset:hover {
      background: #e5e7eb;
    }
    .note {
      margin-top: 2rem;
      color: #9ca3af;
      font-size: 0.8rem;
    }
    .note code {
      background: #f3f4f6;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 0.75rem;
    }
  </style>
</head>
<body>
  <div class="counter">
    <h1>Goscript Counter</h1>
    <p class="subtitle">GoScript 2.0 — reactive, zero JavaScript</p>

    <div id="count">0</div>

    <div class="buttons">
      <button class="dec"
        gs-trigger="click"
        gs-get="/api/counter?n=-1"
        gs-target="#count"
        gs-swap="innerHTML">
        −
      </button>
      <button class="inc"
        gs-trigger="click"
        gs-get="/api/counter?n=1"
        gs-target="#count"
        gs-swap="innerHTML">
        +
      </button>
      <button class="reset"
        gs-trigger="click"
        gs-get="/api/counter?reset=1"
        gs-target="#count"
        gs-swap="innerHTML">
        Reset
      </button>
    </div>

    <p class="note">
      Powered by <code>gs-trigger</code> <code>gs-get</code> <code>gs-target</code> <code>gs-swap</code>
      — the GoScript reactive runtime
    </p>
  </div>
</body>
</html>`, metadata.Render())

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}, []string{"GET"})

	// Counter API — returns an HTML fragment that goscript swaps into #count
	count := 0
	http.HandleFunc("/api/counter", func(w http.ResponseWriter, r *http.Request) {
		n := r.URL.Query().Get("n")
		_, reset := r.URL.Query()["reset"]

		if reset {
			count = 0
		} else if n != "" {
			delta := 0
			fmt.Sscanf(n, "%d", &delta)
			count += delta
		}

		// Return HTML fragment — goscript swaps this into #count
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "%d", count)
	})

	// Health check
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"%s","counter":%d}`, goscript.Version, count)
	})

	fmt.Printf("  ╔══════════════════════════════╗\n")
	fmt.Printf("  ║  Goscript Counter             ║\n")
	fmt.Printf("  ║  http://localhost:%-11s ║\n", port)
	fmt.Printf("  ╚══════════════════════════════╝\n")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
