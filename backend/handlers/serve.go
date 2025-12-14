package handlers

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

	path := r.URL.Query().Get("path")
	fullPath, err := utils.GetSafePath(path)
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
	// Get path from query
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	// Resolve safe full path
	fullPath, err := utils.GetSafePath(path)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Stat the file
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else if os.IsPermission(err) {
			http.Error(w, "Permission denied", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to access file", http.StatusInternalServerError)
		}
		return
	}

	// Check if directory
	if info.IsDir() {
		http.Error(w, "Cannot serve a directory", http.StatusBadRequest)
		return
	}

	// Determine MIME type
	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Set headers
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	// Inline for viewable types, download for others
	switch mimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp",
		"application/pdf", "text/plain", "text/html", "text/css", "application/javascript":
		w.Header().Set("Content-Disposition", "inline; filename=\""+info.Name()+"\"")
	default:
		w.Header().Set("Content-Disposition", "attachment; filename=\""+info.Name()+"\"")
	}

	// Open file safely
	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsPermission(err) {
			http.Error(w, "Permission denied", http.StatusForbidden)
		} else {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Serve content
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

// Download file
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fullPath, err := utils.GetSafePath(path)
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

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", info.Name()))
	http.ServeFile(w, r, fullPath)
}
