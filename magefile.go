//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target when running "mage" without arguments
var Default = Build

// Generate runs the validation code generator
func Generate() error {
	fmt.Println("Generating validation rules...")
	return sh.RunV("go", "run", "./cmd/genvalidation")
}

// npmInstall runs "npm ci" only when package-lock.json is newer than node_modules.
// This ensures local and CI environments always use the same dependency versions
// without paying the full install cost on every invocation.
func npmInstall() error {
	lockFile := filepath.Join("frontend", "package-lock.json")
	installedMark := filepath.Join("frontend", "node_modules", ".package-lock.json")

	lockInfo, err := os.Stat(lockFile)
	if err != nil {
		return fmt.Errorf("cannot stat %s: %w", lockFile, err)
	}

	if markInfo, err := os.Stat(installedMark); err == nil && markInfo.ModTime().After(lockInfo.ModTime()) {
		return nil // node_modules are up to date
	}

	fmt.Println("Installing frontend dependencies...")
	return sh.RunV("npm", "--prefix", "frontend", "ci")
}

// Frontend builds the React frontend with Vite
func Frontend() error {
	mg.Deps(Generate)

	if err := npmInstall(); err != nil {
		return err
	}

	fmt.Println("Building frontend...")
	return sh.RunV("npm", "--prefix", "frontend", "run", "build")
}

// Backend builds the Go backend binary
func Backend() error {
	mg.Deps(Frontend)

	fmt.Println("Building backend...")
	binaryName := "notebook"
	if runtime.GOOS == "windows" {
		binaryName = "notebook.exe"
	}
	return sh.RunV("go", "build", "-o", binaryName, "./cmd/notebook")
}

// Build builds everything (frontend + backend)
func Build() error {
	return Backend()
}

// Test runs all Go tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "./...")
}

// FrontendLint runs the ESLint TypeScript/React linter for the frontend
func FrontendLint() error {
	if err := npmInstall(); err != nil {
		return err
	}
	fmt.Println("Running frontend linter...")
	return sh.RunV("npm", "--prefix", "frontend", "run", "lint")
}

// Lint runs all linters (frontend ESLint + Go golangci-lint)
func Lint() error {
	mg.Deps(FrontendLint)
	fmt.Println("Running Go linter...")
	return sh.RunV("golangci-lint", "run", "--timeout=5m", "./...")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	artifacts := []string{
		"notebook.exe",      // Windows binary
		"notebook",          // Linux/macOS binary
		"genvalidation.exe", // Windows code generator binary
		"genvalidation",     // Linux/macOS code generator binary
		filepath.Join("internal", "web", "frontend", "dist"),
		filepath.Join("frontend", "dist"),
		filepath.Join("frontend", "src", "generated"),
		"tsnet-state",
	}

	for _, artifact := range artifacts {
		if err := os.RemoveAll(artifact); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove %s: %v\n", artifact, err)
		}
	}

	fmt.Println("✓ Cleanup complete")
	return nil
}

// Dev prints development mode instructions
func Dev() {
	fmt.Println("To run in development mode, open two terminals:")
	fmt.Println("")
	fmt.Println("Terminal 1: cd frontend && npm run dev")
	fmt.Println("Terminal 2: go run ./cmd/notebook --dev-listen :8080 --verbose")
	fmt.Println("")
	fmt.Println("Then open http://localhost:5173 in your browser")
}

// Install installs Go dependencies
func Install() error {
	fmt.Println("Downloading Go dependencies...")
	return sh.RunV("go", "mod", "download")
}

// Verify runs all verification steps (lint + test)
func Verify() error {
	mg.Deps(Lint, Test)
	fmt.Println("✓ All verification checks passed")
	return nil
}
