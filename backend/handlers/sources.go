package handlers

import (
	"net/http"

	"filemanager/config"
	"filemanager/utils"
)

func GetSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.SendJSON(w, http.StatusMethodNotAllowed, utils.Response{Success: false, Message: "Method not allowed"})
		return
	}

	sources := config.GetEnabledSources()
	utils.SendJSON(w, http.StatusOK, utils.Response{Success: true, Data: sources})
}
