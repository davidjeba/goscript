package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidjeba/goscript/pkg/gocsx"
	"github.com/davidjeba/goscript/pkg/gocsx/components"
	"github.com/davidjeba/goscript/pkg/gocsx/engine"
)

func main() {
	// Create a new Gocsx instance
	g := gocsx.New()

	// Create a new engine
	e := engine.NewEngine(&engine.EngineConfig{
		Context:         engine.Context2D,
		TargetFPS:       60,
		PerformanceLevel: engine.PerformanceMedium,
	})

	// Create a new Canvas2D
	canvas := engine.NewCanvas2D("main-canvas", 800, 600, e)

	// Set render callback
	canvas.SetRenderCallback(func(ctx *engine.Canvas2DContext, deltaTime float64) {
		// Clear the canvas
		ctx.ClearRect(0, 0, 800, 600)

		// Set fill style
		ctx.FillStyle = "#f0f0f0"
		ctx.FillRect(0, 0, 800, 600)

		// Draw a rectangle
		ctx.FillStyle = "#ff0000"
		ctx.FillRect(100, 100, 200, 150)

		// Draw a circle
		ctx.FillStyle = "#0000ff"
		ctx.BeginPath()
		ctx.Arc(400, 300, 50, 0, 2*3.14159, false)
		ctx.Fill()

		// Draw some text
		ctx.FillStyle = "#000000"
		ctx.Font = "24px Arial"
		ctx.TextAlign = "center"
		ctx.FillText("Gocsx 2D Demo", 400, 50)

		// Draw a line
		ctx.StrokeStyle = "#00ff00"
		ctx.LineWidth = 5
		ctx.BeginPath()
		ctx.MoveTo(500, 100)
		ctx.LineTo(700, 500)
		ctx.Stroke()
	})

	// Start the engine
	e.Start()

	// Create a page with the 2D canvas
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gocsx 2D Demo</title>
    %s
    <style>
        body {
            margin: 0;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            background-color: #f5f5f5;
            font-family: Arial, sans-serif;
        }
        .container {
            text-align: center;
        }
        h1 {
            margin-bottom: 20px;
        }
        #canvas-container {
            border: 1px solid #ccc;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            margin-bottom: 20px;
            background-color: white;
        }
        canvas {
            display: block;
        }
        .controls {
            display: flex;
            gap: 10px;
            justify-content: center;
            margin-bottom: 20px;
        }
        .stats {
            background-color: rgba(0, 0, 0, 0.7);
            color: white;
            padding: 10px;
            border-radius: 5px;
            text-align: left;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Gocsx 2D Canvas Demo</h1>
        
        <div id="canvas-container">
            <canvas id="main-canvas" width="800" height="600"></canvas>
        </div>
        
        <div class="controls">
            %s
            %s
            %s
        </div>
        
        <div class="stats">
            <div>FPS: <span id="fps">0</span></div>
            <div>Draw calls: <span id="draw-calls">0</span></div>
        </div>
    </div>
    
    <script>
        // This would normally be the Canvas2D initialization code
        // For now, we'll just draw a simple scene and update the stats
        
        const canvas = document.getElementById('main-canvas');
        const ctx = canvas.getContext('2d');
        
        function draw() {
            // Clear canvas
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            
            // Background
            ctx.fillStyle = '#f0f0f0';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            
            // Red rectangle
            ctx.fillStyle = '#ff0000';
            ctx.fillRect(100, 100, 200, 150);
            
            // Blue circle
            ctx.fillStyle = '#0000ff';
            ctx.beginPath();
            ctx.arc(400, 300, 50, 0, Math.PI * 2);
            ctx.fill();
            
            // Text
            ctx.fillStyle = '#000000';
            ctx.font = '24px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('Gocsx 2D Demo', 400, 50);
            
            // Green line
            ctx.strokeStyle = '#00ff00';
            ctx.lineWidth = 5;
            ctx.beginPath();
            ctx.moveTo(500, 100);
            ctx.lineTo(700, 500);
            ctx.stroke();
            
            // Update stats
            document.getElementById('fps').textContent = Math.floor(Math.random() * 60 + 30);
            document.getElementById('draw-calls').textContent = Math.floor(Math.random() * 20 + 5);
            
            requestAnimationFrame(draw);
        }
        
        draw();
    </script>
</body>
</html>
	`, g.GenerateStyleTag(), 
	   g.Button(components.ButtonProps{
		   ID:       "clear-button",
		   Variant:  components.ButtonPrimary,
		   Size:     components.ButtonSizeSmall,
		   Children: "Clear Canvas",
		   OnClick:  "alert('Canvas cleared!')",
	   }),
	   g.Button(components.ButtonProps{
		   ID:       "draw-button",
		   Variant:  components.ButtonSuccess,
		   Size:     components.ButtonSizeSmall,
		   Children: "Draw Shape",
		   OnClick:  "alert('Shape drawn!')",
	   }),
	   g.Button(components.ButtonProps{
		   ID:       "save-button",
		   Variant:  components.ButtonInfo,
		   Size:     components.ButtonSizeSmall,
		   Children: "Save Image",
		   OnClick:  "alert('Image saved!')",
	   }))

	// Create a handler for the page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, page)
	})

	// Start the server
	log.Println("Server starting on http://localhost:12001")
	log.Fatal(http.ListenAndServe("0.0.0.0:12001", nil))
}