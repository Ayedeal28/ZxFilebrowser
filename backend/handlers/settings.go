package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"filemanager/config"
	"filemanager/utils"
)

type UserSettings struct {
	DarkMode   bool `json:"darkMode"`
	SidebarPin bool `json:"sidebarPin"`
	Encryption bool `json:"encryption"`
}

func getSettingsFile() string {
	sources := config.GetEnabledSources()
	if len(sources) > 0 {
		return filepath.Join(sources[0].Path, ".settings.json")
	}
	return ".settings.json"
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{Success: false, Message: "Method not allowed"})
		return
	}

	settings := UserSettings{DarkMode: false, SidebarPin: true, Encryption: false}

	data, err := os.ReadFile(getSettingsFile())
	if err == nil {
		json.Unmarshal(data, &settings)
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{Success: true, Data: settings})
}

func SaveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{Success: false, Message: "Method not allowed"})
		return
	}

	var settings UserSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		utils.SendJSON(w, http.StatusBadRequest, utils.Response{Success: false, Message: "Invalid request"})
		return
	}

	data, _ := json.MarshalIndent(settings, "", "  ")
	if err := os.WriteFile(getSettingsFile(), data, 0644); err != nil {
		utils.SendJSON(w, http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to save"})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.Response{Success: true, Message: "Settings saved"})
}
