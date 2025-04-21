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
	e := engine.NewEngine(nil)

	// Create a new WebGPU instance
	webgpu := engine.NewWebGPU()

	// Create a new Three.js scene
	scene := engine.NewThreeJSScene(e, webgpu)

	// Create a camera
	scene.CreateCamera("main-camera", "Main Camera", [3]float64{0, 0, 5}, [3]float64{0, 0, 0})

	// Create a light
	scene.CreateLight("main-light", "Main Light", [3]float64{1, 1, 1}, [3]float64{1, 1, 1}, 1.0, "directional")

	// Create some cubes
	scene.CreateCube("cube1", "Cube 1", [3]float64{-1, 0, 0}, 1.0, [3]float64{1, 0, 0})
	scene.CreateCube("cube2", "Cube 2", [3]float64{1, 0, 0}, 1.0, [3]float64{0, 1, 0})
	scene.CreateSphere("sphere1", "Sphere 1", [3]float64{0, 1, 0}, 0.5, [3]float64{0, 0, 1})

	// Set renderer options
	scene.SetSize(800, 600)
	scene.SetClearColor([4]float64{0.1, 0.1, 0.1, 1.0})
	scene.EnableShadows(true)

	// Start the engine
	e.Start()

	// Create a page with the 3D scene
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gocsx 3D Demo</title>
    %s
    <style>
        body {
            margin: 0;
            overflow: hidden;
        }
        #canvas-container {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%%;
            height: 100%%;
        }
        #ui-container {
            position: absolute;
            top: 10px;
            right: 10px;
            padding: 10px;
            background-color: rgba(0, 0, 0, 0.5);
            color: white;
            border-radius: 5px;
        }
        .stats {
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div id="canvas-container">
        <!-- WebGPU canvas would be inserted here -->
        <canvas id="webgpu-canvas" width="800" height="600"></canvas>
    </div>
    
    <div id="ui-container">
        <div class="stats">
            <div>FPS: <span id="fps">0</span></div>
            <div>Draw calls: <span id="draw-calls">0</span></div>
            <div>Triangles: <span id="triangles">0</span></div>
        </div>
        
        <div class="controls">
            %s
        </div>
    </div>
    
    <script>
        // This would normally be the WebGPU initialization code
        // For now, we'll just update the stats
        function updateStats() {
            document.getElementById('fps').textContent = Math.floor(Math.random() * 60 + 30);
            document.getElementById('draw-calls').textContent = Math.floor(Math.random() * 100 + 50);
            document.getElementById('triangles').textContent = Math.floor(Math.random() * 10000 + 5000);
            
            requestAnimationFrame(updateStats);
        }
        
        updateStats();
    </script>
</body>
</html>
	`, g.GenerateStyleTag(), g.Button(components.ButtonProps{
		ID:       "toggle-rotation",
		Variant:  components.ButtonPrimary,
		Size:     components.ButtonSizeSmall,
		Children: "Toggle Rotation",
		OnClick:  "alert('Rotation toggled!')",
	}))

	// Create a handler for the page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, page)
	})

	// Start the server
	log.Println("Server starting on http://localhost:12000")
	log.Fatal(http.ListenAndServe("0.0.0.0:12000", nil))
}