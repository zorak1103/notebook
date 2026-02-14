package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zorak1103/notebook/internal/db"
	"github.com/zorak1103/notebook/internal/tsapp"
	"github.com/zorak1103/notebook/internal/web"
)

func main() {
	// Parse command-line flags
	var (
		devListen = flag.String("dev-listen", "", "Development mode: listen on this address (e.g., :8080) without Tailscale")
		hostname  = flag.String("hostname", "notebook", "Tailscale hostname for the service")
		stateDir  = flag.String("state-dir", "tsnet-state", "Tailscale state directory")
		dbPath    = flag.String("db", "notebook.db", "SQLite database file path")
		verbose   = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Open database
	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	defer database.Close()

	// Determine if running in dev mode
	devMode := *devListen != ""

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the application
	tsApp, listener := setupListener(ctx, devMode, *devListen, *hostname, *stateDir)
	defer closeTsApp(tsApp)

	// Create and start HTTP server
	httpServer := createHTTPServer(tsApp, database, devMode, *verbose)
	startServer(httpServer, listener)
}

func setupListener(ctx context.Context, devMode bool, devListen, hostname, stateDir string) (*tsapp.App, net.Listener) {
	if devMode {
		fmt.Printf("Starting in development mode on %s\n", devListen)
		lc := &net.ListenConfig{}
		listener, err := lc.Listen(ctx, "tcp", devListen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		return nil, listener
	}

	return setupTailscale(ctx, hostname, stateDir)
}

func setupTailscale(ctx context.Context, hostname, stateDir string) (*tsapp.App, net.Listener) {
	fmt.Printf("Starting Tailscale service as '%s'\n", hostname)
	tsApp := tsapp.New(hostname, stateDir)

	if err := tsApp.Up(ctx); err != nil {
		log.Fatalf("failed to start Tailscale: %v", err)
	}

	listener, err := tsApp.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("failed to create Tailscale listener: %v", err)
	}

	return tsApp, listener
}

func createHTTPServer(tsApp *tsapp.App, database *db.DB, devMode, verbose bool) *http.Server {
	webServer := web.NewServer(tsApp, database, devMode, verbose)
	return &http.Server{
		Handler:      webServer.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func startServer(httpServer *http.Server, listener net.Listener) {
	serverErrors := make(chan error, 1)
	go func() {
		fmt.Printf("HTTP server listening on %s\n", listener.Addr())
		serverErrors <- httpServer.Serve(listener)
	}()

	handleShutdown(httpServer, serverErrors)
}

func handleShutdown(httpServer *http.Server, serverErrors chan error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal %v, shutting down gracefully...\n", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("error during shutdown: %v", err)
		}

		fmt.Println("Shutdown complete")
	}
}

func closeTsApp(tsApp *tsapp.App) {
	if tsApp != nil {
		_ = tsApp.Close()
	}
}
