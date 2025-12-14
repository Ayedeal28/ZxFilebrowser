.PHONY: dev build run clean

dev-frontend:
    cd frontend && npm run dev

dev-backend:
    cd backend && go run main.go

build-frontend:
    cd frontend && npm run build
    rm -rf backend/dist
    cp -r frontend/dist backend/dist

build-backend:
    cd backend && go build -o ../bin/filebrowser

build: build-frontend build-backend

run:
    ./bin/filebrowser

clean:
    rm -rf backend/dist
    rm -rf frontend/dist
    rm -rf bin
