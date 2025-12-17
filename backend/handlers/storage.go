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
	sourceID := r.URL.Query().Get("source")
	if sourceID == "" {
		sources := config.GetEnabledSources()
		if len(sources) > 0 {
			sourceID = sources[0].ID
		}
	}

	var source *config.Source
	for _, s := range config.AppConfig.Sources {
		if s.ID == sourceID {
			source = &s
			break
		}
	}

	if source == nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{Success: false, Message: "Source not found"})
		return
	}

	info, err := getDiskUsage(source.Path)
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
