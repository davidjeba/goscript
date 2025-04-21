package main

import (
        "fmt"
        "log"
        "net/http"
        "os"
        "path/filepath"

        "github.com/davidjeba/goscript/pkg/components"
)

func main() {
        // Create the home page
        home := components.NewGoUIXHomePage("home", nil)

        // Create HTML template
        htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoUIX Demo</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f9f9f9;
        }
        
        .draggable {
            cursor: move;
            position: absolute;
            z-index: 100;
            box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
            transition: box-shadow 0.2s;
        }
        
        .draggable:hover {
            box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3);
        }
        
        button:hover {
            opacity: 0.9;
            transform: translateY(-1px);
        }
        
        button:active {
            transform: translateY(1px);
        }
    </style>
</head>
<body>
    %s
</body>
</html>`

        // Create HTTP server
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // Render the home page
                html := home.Render()
                
                // Insert into template
                fullHTML := fmt.Sprintf(htmlTemplate, html)
                
                // Set content type
                w.Header().Set("Content-Type", "text/html")
                
                // Write response
                w.Write([]byte(fullHTML))
        })

        // Handle static files
        http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
                // Get file path
                path := r.URL.Path[len("/static/"):]
                
                // Get current directory
                dir, err := os.Getwd()
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                
                // Create file path
                filePath := filepath.Join(dir, "static", path)
                
                // Check if file exists
                if _, err := os.Stat(filePath); os.IsNotExist(err) {
                        http.NotFound(w, r)
                        return
                }
                
                // Serve file
                http.ServeFile(w, r, filePath)
        })

        // Handle API requests
        http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
                // Set content type
                w.Header().Set("Content-Type", "application/json")
                
                // Get path
                path := r.URL.Path[len("/api/"):]
                
                // Handle different API endpoints
                if path == "counter/increment" {
                        // Get counter ID
                        counterID := r.URL.Query().Get("id")
                        if counterID == "" {
                                http.Error(w, "Missing counter ID", http.StatusBadRequest)
                                return
                        }
                        
                        // Return success
                        w.Write([]byte(`{"success": true}`))
                        return
                }
                
                // Handle unknown endpoint
                http.NotFound(w, r)
        })

        // Start server
        port := os.Getenv("PORT")
        if port == "" {
                port = "12000"
        }
        
        // Create static directory if it doesn't exist
        if _, err := os.Stat("static"); os.IsNotExist(err) {
                os.Mkdir("static", 0755)
        }
        
        // Log server start
        log.Printf("Server starting on http://localhost:%s", port)
        
        // Start server
        log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}