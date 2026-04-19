package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidjeba/goscript/pkg/gocsx"
	"github.com/davidjeba/goscript/pkg/gocsx/components"
)

func main() {
	// Create a new Gocsx instance
	g := gocsx.New()

	// Create a button
	button := g.Button(components.ButtonProps{
		ID:       "my-button",
		Variant:  components.ButtonPrimary,
		Size:     components.ButtonSizeLarge,
		Children: "Click Me",
		OnClick:  "alert('Hello, Gocsx!')",
	})

	// Create a card
	card := g.Card(components.CardProps{
		ID:        "my-card",
		Title:     "Gocsx Card",
		Subtitle:  "A powerful CSS framework for Go",
		Body:      "This is a card component built with Gocsx, a CSS framework for Go that can be used to build web, mobile, and AR/VR applications.",
		Footer:    "Footer content",
		ClassName: "shadow",
	})

	// Create a page
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gocsx Demo</title>
    %s
</head>
<body>
    <div class="container py-5">
        <h1 class="text-center mb-5">Gocsx Demo</h1>
        
        <div class="row mb-5">
            <div class="col-md-6 offset-md-3">
                <div class="text-center">
                    %s
                </div>
            </div>
        </div>
        
        <div class="row">
            <div class="col-md-6 offset-md-3">
                %s
            </div>
        </div>
    </div>
</body>
</html>
	`, g.GenerateStyleTag(), button, card)

	// Create a handler for the page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, page)
	})

	// Start the server
	log.Println("Server starting on http://localhost:12000")
	log.Fatal(http.ListenAndServe("0.0.0.0:12000", nil))
}