package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed dist/*
var frontendFS embed.FS

func serveFrontend() http.Handler {
	// Try to serve embedded frontend
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		// If dist doesn't exist (development mode), just return a handler that shows a message
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<html>
					<body>
						<h1>Development Mode</h1>
						<p>Frontend is running separately on <a href="http://localhost:5173">http://localhost:5173</a></p>
						<p>API is ready at <a href="/api">/api</a></p>
					</body>
				</html>
			`))
		})
	}
	return http.FileServer(http.FS(distFS))
}

// FileInfo represents file/folder metadata
type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Ext     string    `json:"ext"`
}

// Response structure for API responses
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Config holds server configuration
type Config struct {
	RootDir string
	Port    string
}

var config Config

func main() {
	// Configuration
	config = Config{
		RootDir: "C:\\Users\\ayede\\Desktop\\TestFile",
		Port:    getEnv("PORT", ":8080"),
	}

	// Create root directory if it doesn't exist
	if err := os.MkdirAll(config.RootDir, 0755); err != nil {
		log.Fatal("Failed to create root directory:", err)
	}

	// API routes
	http.HandleFunc("/api/list", corsMiddleware(listHandler))
	http.HandleFunc("/api/create", corsMiddleware(createHandler))
	http.HandleFunc("/api/delete", corsMiddleware(deleteHandler))
	http.HandleFunc("/api/copy", corsMiddleware(copyHandler))
	http.HandleFunc("/api/move", corsMiddleware(moveHandler))
	http.HandleFunc("/api/rename", corsMiddleware(renameHandler))
	http.HandleFunc("/api/preview", corsMiddleware(previewHandler))
	http.HandleFunc("/api/serve", corsMiddleware(serveFileHandler))
	http.HandleFunc("/api/upload", corsMiddleware(uploadHandler))
	http.HandleFunc("/api/download", corsMiddleware(downloadHandler))
	http.HandleFunc("/api/info", corsMiddleware(infoHandler))

	http.Handle("/", serveFrontend())

	log.Printf("🚀 Server starting on http://localhost%s", config.Port)
	log.Printf("📁 Serving files from: %s", config.RootDir)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Helper function to get safe path
func getSafePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	fullPath := filepath.Join(config.RootDir, cleanPath)

	absRoot, err := filepath.Abs(config.RootDir)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, absRoot) {
		return "", fmt.Errorf("invalid path: outside root directory")
	}

	return absPath, nil
}

// Send JSON response
func sendJSON(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Get file info
func infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	path := r.URL.Query().Get("path")
	fullPath, err := getSafePath(path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		sendJSON(w, http.StatusNotFound, Response{Success: false, Message: "File not found"})
		return
	}

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Ext:     filepath.Ext(info.Name()),
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Data: fileInfo})
}

// List files and folders
func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	fullPath, err := getSafePath(path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to read directory"})
		return
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		relativePath := filepath.Join(path, entry.Name())
		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    relativePath,
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Ext:     filepath.Ext(entry.Name()),
		})
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Data: files})
}

// Create file or folder
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Path  string `json:"path"`
		IsDir bool   `json:"isDir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid request"})
		return
	}

	fullPath, err := getSafePath(req.Path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	if req.IsDir {
		err = os.MkdirAll(fullPath, 0755)
	} else {
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to create parent directory"})
			return
		}
		_, err = os.Create(fullPath)
	}

	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to create"})
		return
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Message: "Created successfully"})
}

// Delete file or folder
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid request"})
		return
	}

	fullPath, err := getSafePath(req.Path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	if err := os.RemoveAll(fullPath); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to delete"})
		return
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Message: "Deleted successfully"})
}

// Copy file or folder
func copyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid request"})
		return
	}

	srcPath, err := getSafePath(req.Source)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid source"})
		return
	}

	dstPath, err := getSafePath(req.Destination)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid destination"})
		return
	}

	if err := copyPath(srcPath, dstPath); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to copy"})
		return
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Message: "Copied successfully"})
}

// Move file or folder
func moveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid request"})
		return
	}

	srcPath, err := getSafePath(req.Source)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid source"})
		return
	}

	dstPath, err := getSafePath(req.Destination)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid destination"})
		return
	}

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to create destination"})
		return
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to move"})
		return
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Message: "Moved successfully"})
}

// Rename file or folder
func renameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	var req struct {
		Path    string `json:"path"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Invalid request"})
		return
	}

	oldPath, err := getSafePath(req.Path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to rename"})
		return
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Message: "Renamed successfully"})
}

// Preview file
func previewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	path := r.URL.Query().Get("path")
	fullPath, err := getSafePath(path)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		sendJSON(w, http.StatusNotFound, Response{Success: false, Message: "File not found"})
		return
	}

	if info.IsDir() {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Cannot preview directory"})
		return
	}

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   false,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Ext:     filepath.Ext(info.Name()),
	}

	// Read content for small text files
	if info.Size() < 1024*1024 {
		content, err := os.ReadFile(fullPath)
		if err == nil && isTextFile(content) {
			sendJSON(w, http.StatusOK, Response{
				Success: true,
				Data: map[string]interface{}{
					"info":    fileInfo,
					"content": string(content),
				},
			})
			return
		}
	}

	sendJSON(w, http.StatusOK, Response{Success: true, Data: map[string]interface{}{"info": fileInfo}})
}

// Serve file
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fullPath, err := getSafePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, fullPath)
}

// Download file
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fullPath, err := getSafePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", info.Name()))
	http.ServeFile(w, r, fullPath)
}

// Upload file
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	if err := r.ParseMultipartForm(100 << 20); err != nil { // 100MB max
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "Failed to parse form"})
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: "No file provided"})
		return
	}
	defer file.Close()

	destPath := filepath.Join(path, header.Filename)
	fullPath, err := getSafePath(destPath)
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{Success: false, Message: err.Error()})
		return
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to create directory"})
		return
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to create file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{Success: false, Message: "Failed to save file"})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Uploaded successfully",
		Data:    map[string]string{"path": destPath, "name": header.Filename},
	})
}

// Helper: Copy file or directory
func copyPath(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, _ := os.Stat(src)
	return os.Chmod(dst, srcInfo.Mode())
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper: Check if content is text
func isTextFile(content []byte) bool {
	if len(content) == 0 {
		return true
	}

	for i := 0; i < len(content) && i < 512; i++ {
		if content[i] == 0 {
			return false
		}
	}
	return true
}
