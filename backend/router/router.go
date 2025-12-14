package router

import (
	"net/http"

	"filemanager/handlers"
	"filemanager/utils"
)

func SetupRoutes() {
	// Directory operations
	http.HandleFunc("/api/list", utils.CORSMiddleware(handlers.ListDirectory))
	http.HandleFunc("/api/info", utils.CORSMiddleware(handlers.GetInfo))

	// File operations
	http.HandleFunc("/api/create", utils.CORSMiddleware(handlers.CreateItem))
	http.HandleFunc("/api/delete", utils.CORSMiddleware(handlers.DeleteItem))
	http.HandleFunc("/api/rename", utils.CORSMiddleware(handlers.RenameItem))
	http.HandleFunc("/api/copy", utils.CORSMiddleware(handlers.CopyItem))
	http.HandleFunc("/api/move", utils.CORSMiddleware(handlers.MoveItem))
	http.HandleFunc("/api/upload", utils.CORSMiddleware(handlers.UploadFile))

	// File serving
	http.HandleFunc("/api/preview", utils.CORSMiddleware(handlers.PreviewFile))
	http.HandleFunc("/api/serve", utils.CORSMiddleware(handlers.ServeFile))
	http.HandleFunc("/api/download", utils.CORSMiddleware(handlers.DownloadFile))

	// Storage info
	http.HandleFunc("/api/storage", utils.CORSMiddleware(handlers.GetStorageInfo))
}
