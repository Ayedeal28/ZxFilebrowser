package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"filemanager/config"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func GetSafePath(sourceID, path string) (string, error) {
	// Find the source
	var source *config.Source
	for _, s := range config.AppConfig.Sources {
		if s.ID == sourceID && s.Enabled {
			source = &s
			break
		}
	}

	if source == nil {
		return "", fmt.Errorf("source not found or disabled: %s", sourceID)
	}

	cleanPath := filepath.Clean(path)
	fullPath := filepath.Join(source.Path, cleanPath)

	absRoot, err := filepath.Abs(source.Path)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, absRoot) {
		return "", fmt.Errorf("invalid path: outside source directory")
	}

	return absPath, nil
}

func SendJSON(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

	// Log response
	statusText := "SUCCESS"
	if status >= 400 {
		statusText = "ERROR"
	}
	log.Printf("[%s] %d - %s", statusText, status, response.Message)
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func IsTextFile(content []byte) bool {
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
