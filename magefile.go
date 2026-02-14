//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

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

// Frontend builds the React frontend with Vite
func Frontend() error {
	mg.Deps(Generate)

	fmt.Println("Installing frontend dependencies...")
	if err := sh.RunV("npm", "--prefix", "frontend", "ci"); err != nil {
		return err
	}

	fmt.Println("Building frontend...")
	return sh.RunV("npm", "--prefix", "frontend", "run", "build")
}

// Backend builds the Go backend binary
func Backend() error {
	mg.Deps(Frontend)

	fmt.Println("Building backend...")
	return sh.RunV("go", "build", "-o", "notebook.exe", "./cmd/notebook")
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

// Lint runs golangci-lint with timeout
func Lint() error {
	fmt.Println("Running linter...")
	return sh.RunV("golangci-lint", "run", "--timeout=5m", "./...")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	artifacts := []string{
		"notebook.exe",
		"notebook",
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
