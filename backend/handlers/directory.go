package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"filemanager/utils"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	Ext     string    `json:"ext"`
}

func ListDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	fullPath, err := utils.GetSafePath(path)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to read directory",
		})
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

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Data:    files,
	})
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
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

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    path,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		Ext:     filepath.Ext(info.Name()),
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Data:    fileInfo,
	})
}
