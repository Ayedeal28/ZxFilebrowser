package main

import (
	"embed"
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
	if err := os.MkdirAll(config.AppConfig.RootDir, 0755); err != nil {
		log.Fatal("Failed to create root directory:", err)
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

	log.Printf("🚀 Server starting on http://localhost%s", config.AppConfig.Port)
	log.Printf("📁 Serving files from: %s", config.AppConfig.RootDir)
	log.Fatal(http.ListenAndServe(config.AppConfig.Port, nil))
}
