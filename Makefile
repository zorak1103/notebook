.PHONY: all build frontend backend dev clean test

# Build everything
all: build

# Complete build process (frontend first, then backend)
build: frontend backend

# Build the frontend with Vite
frontend:
	cd frontend && npm ci && npm run build

# Build the Go backend
backend:
	go build -o notebook.exe ./cmd/notebook/

# Run tests
test:
	go test ./...

# Development mode instructions
dev:
	@echo "To run in development mode, open two terminals:"
	@echo ""
	@echo "Terminal 1: cd frontend && npm run dev"
	@echo "Terminal 2: go run ./cmd/notebook --dev-listen :8080 --verbose"
	@echo ""
	@echo "Then open http://localhost:5173 in your browser"

# Clean build artifacts
clean:
	rm -f notebook.exe
	rm -rf internal/web/frontend/dist/*
	rm -rf frontend/dist
	rm -rf tsnet-state
