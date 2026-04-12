package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidjeba/goscript/pkg/goscale/api"
	"github.com/davidjeba/goscript/pkg/goscale/db"
	"github.com/davidjeba/goscript/pkg/goscale/edge"
)

func main() {
	// Create a new GoScaleAPI instance
	apiConfig := &api.Config{
		DBConnectionString: "localhost:5432",
		EdgeEnabled:        true,
		EdgeNodes:          []string{"edge-1", "edge-2"},
		CompressionLevel:   5,
		BatchSize:          100,
		Timeout:            time.Second * 30,
		MaxConcurrent:      1000,
		EnableTimeSeries:   true,
		EnableRelationships: true,
		EnableNoCode:       true,
	}
	
	goscaleAPI := api.NewGoScaleAPI(apiConfig)
	
	// Create a schema
	schema := api.NewSchema()
	
	// Add types
	userType := schema.AddType("User", "A user in the system")
	userType.AddField("id", "ID", "The user's ID")
	userType.AddField("name", "String", "The user's name")
	userType.AddField("email", "String", "The user's email")
	userType.AddField("createdAt", "DateTime", "When the user was created")
	
	postType := schema.AddType("Post", "A post in the system")
	postType.AddField("id", "ID", "The post's ID")
	postType.AddField("title", "String", "The post's title")
	postType.AddField("content", "String", "The post's content")
	postType.AddField("author", "User", "The post's author")
	postType.AddField("createdAt", "DateTime", "When the post was created")
	
	// Add queries
	getUserQuery := schema.AddQuery("getUser", "User", "Get a user by ID")
	getUserQuery.AddArg("id", "ID", nil, "The user's ID")
	getUserQuery.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// In a real implementation, this would query the database
		return map[string]interface{}{
			"id":        params["id"],
			"name":      "John Doe",
			"email":     "john@example.com",
			"createdAt": time.Now().Format(time.RFC3339),
		}, nil
	})
	
	getPostsQuery := schema.AddQuery("getPosts", "[Post]", "Get posts by user ID")
	getPostsQuery.AddArg("userId", "ID", nil, "The user's ID")
	getPostsQuery.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// In a real implementation, this would query the database
		return []map[string]interface{}{
			{
				"id":        1,
				"title":     "First Post",
				"content":   "This is my first post",
				"author":    params["userId"],
				"createdAt": time.Now().Format(time.RFC3339),
			},
			{
				"id":        2,
				"title":     "Second Post",
				"content":   "This is my second post",
				"author":    params["userId"],
				"createdAt": time.Now().Format(time.RFC3339),
			},
		}, nil
	})
	
	// Add mutations
	createUserMutation := schema.AddMutation("createUser", "User", "Create a new user")
	createUserMutation.AddArg("name", "String", nil, "The user's name")
	createUserMutation.AddArg("email", "String", nil, "The user's email")
	createUserMutation.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// In a real implementation, this would insert into the database
		return map[string]interface{}{
			"id":        123,
			"name":      params["name"],
			"email":     params["email"],
			"createdAt": time.Now().Format(time.RFC3339),
		}, nil
	})
	
	createPostMutation := schema.AddMutation("createPost", "Post", "Create a new post")
	createPostMutation.AddArg("title", "String", nil, "The post's title")
	createPostMutation.AddArg("content", "String", nil, "The post's content")
	createPostMutation.AddArg("authorId", "ID", nil, "The author's ID")
	createPostMutation.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// In a real implementation, this would insert into the database
		return map[string]interface{}{
			"id":        456,
			"title":     params["title"],
			"content":   params["content"],
			"author":    params["authorId"],
			"createdAt": time.Now().Format(time.RFC3339),
		}, nil
	})
	
	// Add subscriptions
	newPostSubscription := schema.AddSubscription("newPost", "Post", "Subscribe to new posts")
	newPostSubscription.AddArg("userId", "ID", nil, "The user's ID")
	newPostSubscription.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// In a real implementation, this would set up a subscription
		return nil, nil
	})
	
	// Apply the schema to the API
	err := goscaleAPI.ApplySchema(schema)
	if err != nil {
		log.Fatalf("Error applying schema: %v", err)
	}
	
	// Create an edge network
	edgeNetwork := edge.NewEdgeNetwork(goscaleAPI)
	
	// Create edge nodes
	edgeConfig1 := &edge.Config{
		ID:               "edge-1",
		Region:           "us-east",
		Capacity:         1000,
		CacheEnabled:     true,
		CacheTTL:         time.Minute * 5,
		DBConfig:         db.DefaultConfig(),
		SyncInterval:     time.Minute * 15,
		MaxConcurrent:    100,
		CompressionLevel: 5,
	}
	
	edgeNode1 := edge.NewEdgeNode(edgeConfig1, goscaleAPI)
	edgeNetwork.AddNode(edgeNode1)
	
	edgeConfig2 := &edge.Config{
		ID:               "edge-2",
		Region:           "us-west",
		Capacity:         1000,
		CacheEnabled:     true,
		CacheTTL:         time.Minute * 5,
		DBConfig:         db.DefaultConfig(),
		SyncInterval:     time.Minute * 15,
		MaxConcurrent:    100,
		CompressionLevel: 5,
	}
	
	edgeNode2 := edge.NewEdgeNode(edgeConfig2, goscaleAPI)
	edgeNetwork.AddNode(edgeNode2)
	
	// Register handlers for edge nodes
	for path, resolver := range goscaleAPI.GetResolvers() {
		edgeNode1.RegisterHandler(path, resolver)
		edgeNode2.RegisterHandler(path, resolver)
	}
	
	// Create HTTP handlers
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		goscaleAPI.ServeHTTP(w, r)
	})
	
	http.HandleFunc("/edge", func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var request struct {
			Path   string                 `json:"path"`
			Params map[string]interface{} `json:"params"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
		defer cancel()
		
		// Process the request through the edge network
		result, err := edgeNetwork.ProcessRequest(ctx, request.Path, request.Params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Return the result
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": result,
		})
	})
	
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// Get metrics from the API and edge nodes
		apiMetrics := goscaleAPI.GetMetrics()
		edge1Metrics := edgeNode1.GetMetrics()
		edge2Metrics := edgeNode2.GetMetrics()
		
		// Return the metrics
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"api": apiMetrics,
			"edge": map[string]interface{}{
				"edge-1": edge1Metrics,
				"edge-2": edge2Metrics,
			},
		})
	})
	
	// Create a simple UI for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoScale API Demo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1 {
            color: #333;
        }
        .card {
            background-color: #f5f5f5;
            border-radius: 5px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input, textarea, select {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            border: none;
            padding: 10px 15px;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #45a049;
        }
        pre {
            background-color: #f9f9f9;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 10px;
            overflow: auto;
        }
        .tabs {
            display: flex;
            margin-bottom: 20px;
        }
        .tab {
            padding: 10px 15px;
            cursor: pointer;
            border: 1px solid #ddd;
            background-color: #f5f5f5;
            margin-right: 5px;
            border-radius: 4px 4px 0 0;
        }
        .tab.active {
            background-color: #fff;
            border-bottom: 1px solid #fff;
        }
        .tab-content {
            display: none;
            border: 1px solid #ddd;
            padding: 20px;
            border-radius: 0 4px 4px 4px;
        }
        .tab-content.active {
            display: block;
        }
    </style>
</head>
<body>
    <h1>GoScale API Demo</h1>
    
    <div class="tabs">
        <div class="tab active" data-tab="query">Query</div>
        <div class="tab" data-tab="mutation">Mutation</div>
        <div class="tab" data-tab="edge">Edge</div>
        <div class="tab" data-tab="metrics">Metrics</div>
    </div>
    
    <div class="tab-content active" id="query-tab">
        <div class="card">
            <h2>Query</h2>
            <div class="form-group">
                <label for="query-type">Query Type</label>
                <select id="query-type">
                    <option value="getUser">Get User</option>
                    <option value="getPosts">Get Posts</option>
                </select>
            </div>
            <div class="form-group">
                <label for="query-params">Parameters (JSON)</label>
                <textarea id="query-params" rows="5">{"id": 123}</textarea>
            </div>
            <button id="run-query">Run Query</button>
        </div>
        <h3>Result</h3>
        <pre id="query-result"></pre>
    </div>
    
    <div class="tab-content" id="mutation-tab">
        <div class="card">
            <h2>Mutation</h2>
            <div class="form-group">
                <label for="mutation-type">Mutation Type</label>
                <select id="mutation-type">
                    <option value="createUser">Create User</option>
                    <option value="createPost">Create Post</option>
                </select>
            </div>
            <div class="form-group">
                <label for="mutation-params">Parameters (JSON)</label>
                <textarea id="mutation-params" rows="5">{"name": "John Doe", "email": "john@example.com"}</textarea>
            </div>
            <button id="run-mutation">Run Mutation</button>
        </div>
        <h3>Result</h3>
        <pre id="mutation-result"></pre>
    </div>
    
    <div class="tab-content" id="edge-tab">
        <div class="card">
            <h2>Edge Computing</h2>
            <div class="form-group">
                <label for="edge-path">Path</label>
                <input type="text" id="edge-path" value="query:getUser">
            </div>
            <div class="form-group">
                <label for="edge-params">Parameters (JSON)</label>
                <textarea id="edge-params" rows="5">{"id": 123}</textarea>
            </div>
            <button id="run-edge">Run Edge Request</button>
        </div>
        <h3>Result</h3>
        <pre id="edge-result"></pre>
    </div>
    
    <div class="tab-content" id="metrics-tab">
        <div class="card">
            <h2>Metrics</h2>
            <button id="get-metrics">Get Metrics</button>
        </div>
        <h3>Result</h3>
        <pre id="metrics-result"></pre>
    </div>
    
    <script>
        // Tab switching
        document.querySelectorAll('.tab').forEach(tab => {
            tab.addEventListener('click', () => {
                document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
                document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
                
                tab.classList.add('active');
                document.getElementById(tab.dataset.tab + '-tab').classList.add('active');
            });
        });
        
        // Query
        document.getElementById('run-query').addEventListener('click', async () => {
            const queryType = document.getElementById('query-type').value;
            const params = JSON.parse(document.getElementById('query-params').value);
            
            try {
                const response = await fetch('/api', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        query: queryType,
                        variables: params,
                        operation: 'query:' + queryType
                    })
                });
                
                const result = await response.json();
                document.getElementById('query-result').textContent = JSON.stringify(result, null, 2);
            } catch (error) {
                document.getElementById('query-result').textContent = 'Error: ' + error.message;
            }
        });
        
        // Mutation
        document.getElementById('run-mutation').addEventListener('click', async () => {
            const mutationType = document.getElementById('mutation-type').value;
            const params = JSON.parse(document.getElementById('mutation-params').value);
            
            try {
                const response = await fetch('/api', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        query: mutationType,
                        variables: params,
                        operation: 'mutation:' + mutationType
                    })
                });
                
                const result = await response.json();
                document.getElementById('mutation-result').textContent = JSON.stringify(result, null, 2);
            } catch (error) {
                document.getElementById('mutation-result').textContent = 'Error: ' + error.message;
            }
        });
        
        // Edge
        document.getElementById('run-edge').addEventListener('click', async () => {
            const path = document.getElementById('edge-path').value;
            const params = JSON.parse(document.getElementById('edge-params').value);
            
            try {
                const response = await fetch('/edge', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        path: path,
                        params: params
                    })
                });
                
                const result = await response.json();
                document.getElementById('edge-result').textContent = JSON.stringify(result, null, 2);
            } catch (error) {
                document.getElementById('edge-result').textContent = 'Error: ' + error.message;
            }
        });
        
        // Metrics
        document.getElementById('get-metrics').addEventListener('click', async () => {
            try {
                const response = await fetch('/metrics');
                const result = await response.json();
                document.getElementById('metrics-result').textContent = JSON.stringify(result, null, 2);
            } catch (error) {
                document.getElementById('metrics-result').textContent = 'Error: ' + error.message;
            }
        });
        
        // Set default parameters based on selected query/mutation
        document.getElementById('query-type').addEventListener('change', () => {
            const queryType = document.getElementById('query-type').value;
            if (queryType === 'getUser') {
                document.getElementById('query-params').value = '{"id": 123}';
            } else if (queryType === 'getPosts') {
                document.getElementById('query-params').value = '{"userId": 123}';
            }
        });
        
        document.getElementById('mutation-type').addEventListener('change', () => {
            const mutationType = document.getElementById('mutation-type').value;
            if (mutationType === 'createUser') {
                document.getElementById('mutation-params').value = '{"name": "John Doe", "email": "john@example.com"}';
            } else if (mutationType === 'createPost') {
                document.getElementById('mutation-params').value = '{"title": "New Post", "content": "This is a new post", "authorId": 123}';
            }
        });
    </script>
</body>
</html>
        `;
        
        fmt.Fprint(w, html);
    });
	
	// Start the server
	log.Println("Server starting on http://localhost:12001")
	log.Fatal(http.ListenAndServe("0.0.0.0:12001", nil))
}