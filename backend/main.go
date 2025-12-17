package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"filemanager/config"
	"filemanager/router"
)

//go:embed all:dist
var frontendFS embed.FS

func main() {
	// Initialize configuration
	config.Init()

	// Create root directory if it doesn't exist
	for _, source := range config.AppConfig.Sources {
		if source.Enabled && source.Type == "local" {
			if err := os.MkdirAll(source.Path, 0755); err != nil {
				log.Printf("Warning: Failed to create directory for %s: %v", source.Name, err)
			}
		}
	}

	// Setup API routes
	router.SetupRoutes()

	// Serve frontend
	distFS, err := fs.Sub(frontendFS, "dist")
	if err == nil {
		// Production: serve embedded frontend
		http.Handle("/", http.FileServer(http.FS(distFS)))
		log.Println("📦 Serving embedded frontend")
	} else {
		// Development: show API info
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<html>
					<head><title>ZxFileBrowser API</title></head>
					<body style="font-family: sans-serif; padding: 40px; max-width: 800px; margin: 0 auto;">
						<h1>🚀 ZxFileBrowser API Server</h1>
						<p>API is running and ready to accept requests.</p>
						<p><strong>Frontend:</strong> <a href="http://localhost:5173">http://localhost:5173</a> (Development)</p>
						<h2>API Endpoints:</h2>
						<ul>
							<li>GET /api/list - List directory contents</li>
							<li>GET /api/info - Get file/folder info</li>
							<li>GET /api/preview - Preview file</li>
							<li>GET /api/serve - Serve file</li>
							<li>GET /api/download - Download file</li>
							<li>POST /api/create - Create file/folder</li>
							<li>POST /api/upload - Upload file</li>
							<li>POST /api/rename - Rename item</li>
							<li>POST /api/copy - Copy item</li>
							<li>POST /api/move - Move item</li>
							<li>DELETE /api/delete - Delete item</li>
						</ul>
					</body>
				</html>
			`))
		})
		log.Println("⚠️  Running in development mode (frontend not embedded)")
	}

	port := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("🚀 Server starting on http://localhost%s", port)
	for _, src := range config.GetEnabledSources() {
		log.Printf("📁 Serving: %s from %s", src.Name, src.Path)
	}
	log.Fatal(http.ListenAndServe(port, nil))
}
