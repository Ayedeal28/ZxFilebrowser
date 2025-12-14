.PHONY: dev build run clean install build-and-run

# Install dependencies
install:
	cd frontend && npm install
	cd backend && go mod tidy

# Development mode
dev-frontend:
	cd frontend && npm run dev

dev-backend:
	cd backend && go run main.go

# Build frontend
build-frontend:
	cd frontend && npm run build
	@if exist backend\dist rmdir /s /q backend\dist
	@xcopy /e /i /y frontend\dist backend\dist

# Build backend
build-backend:
	cd backend && go build -o zxfilebrowser.exe

# Build everything
build: build-frontend build-backend
	@echo ✅ Build complete! Binary at .\backend\zxfilebrowser.exe

# Build and run automatically
build-and-run: build
	@echo 🚀 Starting ZxFileBrowser...
	.\backend\zxfilebrowser.exe

# Run the built binary
run:
	.\backend\zxfilebrowser.exe

# Clean build artifacts
clean:
	@if exist backend\dist rmdir /s /q backend\dist
	@if exist frontend\dist rmdir /s /q frontend\dist
	@if exist backend\zxfilebrowser.exe del backend\zxfilebrowser.exe
	@if exist bin rmdir /s /q bin
	@echo ✅ Cleaned all build artifacts