package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"filemanager/utils"
)

// Create file or folder
func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Path  string `json:"path"`
		IsDir bool   `json:"isDir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	fullPath, err := utils.GetSafePath(req.Path)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if req.IsDir {
		err = os.MkdirAll(fullPath, 0755)
	} else {
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to create parent directory",
			})
			return
		}
		_, err = os.Create(fullPath)
	}

	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Created successfully",
	})
}

// Delete file or folder
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	fullPath, err := utils.GetSafePath(req.Path)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := os.RemoveAll(fullPath); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to delete",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Deleted successfully",
	})
}

// Rename file or folder
func RenameItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Path    string `json:"path"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	oldPath, err := utils.GetSafePath(req.Path)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	newPath := filepath.Join(filepath.Dir(oldPath), req.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to rename",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Renamed successfully",
	})
}

// Copy file or folder
func CopyItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	srcPath, err := utils.GetSafePath(req.Source)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid source",
		})
		return
	}

	dstPath, err := utils.GetSafePath(req.Destination)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid destination",
		})
		return
	}

	if err := copyPath(srcPath, dstPath); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to copy",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Copied successfully",
	})
}

// Move file or folder
func MoveItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	srcPath, err := utils.GetSafePath(req.Source)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid source",
		})
		return
	}

	dstPath, err := utils.GetSafePath(req.Destination)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid destination",
		})
		return
	}

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create destination",
		})
		return
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to move",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Moved successfully",
	})
}

// Upload file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Failed to parse form",
		})
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "No file provided",
		})
		return
	}
	defer file.Close()

	destPath := filepath.Join(path, header.Filename)
	fullPath, err := utils.GetSafePath(destPath)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create directory",
		})
		return
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create file",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{
		Success: true,
		Message: "Uploaded successfully",
		Data: map[string]string{
			"path": destPath,
			"name": header.Filename,
		},
	})
}

// Helper functions for copy operations
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
