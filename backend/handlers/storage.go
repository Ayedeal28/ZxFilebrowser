package handlers

import (
	"net/http"
	"path/filepath"

	"filemanager/config"
	"filemanager/utils"
)

// Storage information structure
type StorageInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Path  string `json:"path"`
}

// GetStorageInfo returns disk usage information
func GetStorageInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	info, err := getDiskUsage(config.AppConfig.RootDir)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to get storage info: " + err.Error(),
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Data:    info,
	})
}

// getDiskUsage gets disk usage for a given path
func getDiskUsage(path string) (*StorageInfo, error) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Get disk stats based on OS
	total, free, err := getDiskStats(absPath)
	if err != nil {
		return nil, err
	}

	used := total - free

	return &StorageInfo{
		Total: total,
		Used:  used,
		Free:  free,
		Path:  absPath,
	}, nil
}
