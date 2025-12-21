package handlers

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"filemanager/utils"
)

// Preview file content
func PreviewFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	sourceID := r.URL.Query().Get("source")
	path := r.URL.Query().Get("path")
	path = strings.ReplaceAll(path, "\\", "/")
	fullPath, err := utils.GetSafePath(sourceID, path)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		utils.SendJSON(w, http.StatusNotFound, utils.Response{
			Success: false,
			Message: "File not found",
		})
		return
	}

	if info.IsDir() {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Cannot preview directory",
		})
		return
	}

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   false,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	// Read content for small text files
	if info.Size() < 1024*1024 {
		content, err := os.ReadFile(fullPath)
		if err == nil && utils.IsTextFile(content) {
			utils.SendJSON(w, http.StatusOK, utils.Response{
				Success: true,
				Data: map[string]interface{}{
					"info":    fileInfo,
					"content": string(content),
				},
			})
			return
		}
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Data: map[string]interface{}{
			"info": fileInfo,
		},
	})
}

// Serve file for viewing
func ServeFile(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source")
	path := r.URL.Query().Get("path")
	path = strings.ReplaceAll(path, "\\", "/")

	if sourceID == "" || path == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	fullPath, err := utils.GetSafePath(sourceID, path)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to access file", http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		http.Error(w, "Cannot serve a directory", http.StatusBadRequest)
		return
	}

	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileName := info.Name()

	setContentDisposition(w, r, fileName)

	// Determine MIME type
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(fileName)))
	if contentType == "" {
		contentType = "text/plain" // fallback for unknown types
	}
	w.Header().Set("Content-Type", contentType)

	// Serve inline instead of forcing download
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", fileName))

	// Prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	log.Printf("Serving: %s", fileName)

	http.ServeContent(w, r, fileName, info.ModTime(), file)
}

// setContentDisposition sets the Content-Disposition header based on file type
func setContentDisposition(w http.ResponseWriter, r *http.Request, fileName string) {
	ext := strings.ToLower(filepath.Ext(fileName))
	mimeType := mime.TypeByExtension(ext)

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimeType)

	// Force download only if explicitly requested
	if r.URL.Query().Get("download") == "true" {
		w.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		return
	}

	// Default: let browser decide (inline if supported)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", fileName))

}

// Download file
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	sourceID := r.URL.Query().Get("source")
	path := r.URL.Query().Get("path")
	path = strings.ReplaceAll(path, "\\", "/")
	fullPath, err := utils.GetSafePath(sourceID, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		http.Error(w, "Cannot download directory", http.StatusBadRequest)
		return
	}
	log.Printf("Downloading: %s (%d bytes)", info.Name(), info.Size())
	// Set proper headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", info.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

	http.ServeFile(w, r, fullPath)
}
