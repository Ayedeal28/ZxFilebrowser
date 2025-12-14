package utils

import (
	"encoding/json"
	"fmt"
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

func GetSafePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	fullPath := filepath.Join(config.AppConfig.RootDir, cleanPath)

	absRoot, err := filepath.Abs(config.AppConfig.RootDir)
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

func SendJSON(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
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
