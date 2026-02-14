package tsapp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"tailscale.com/client/local"
	"tailscale.com/tsnet"
)

// UserInfo represents Tailscale user information returned by WhoIs
type UserInfo struct {
	DisplayName   string `json:"displayName"`
	LoginName     string `json:"loginName"`
	ProfilePicURL string `json:"profilePicURL"`
	NodeName      string `json:"nodeName"`
	NodeID        string `json:"nodeID"`
}

// App wraps a Tailscale tsnet.Server and local client for managing
// the Tailscale network connection and user lookups
type App struct {
	server *tsnet.Server
	lc     *local.Client
}

// New creates a new Tailscale app with the given hostname and state directory
func New(hostname, stateDir string) *App {
	srv := &tsnet.Server{
		Hostname: hostname,
		Dir:      stateDir,
		Logf: func(format string, args ...interface{}) {
			fmt.Printf("[tsnet] "+format+"\n", args...)
		},
	}

	return &App{
		server: srv,
		lc:     nil, // Will be initialized in Up()
	}
}

// Up starts the Tailscale connection and waits for it to be ready
func (a *App) Up(ctx context.Context) error {
	status, err := a.server.Up(ctx)
	if err != nil {
		return fmt.Errorf("tailscale up: %w", err)
	}

	// Get the local client from the server
	a.lc, err = a.server.LocalClient()
	if err != nil {
		return fmt.Errorf("failed to get local client: %w", err)
	}

	fmt.Printf("Tailscale connected. Node: %s\n", status.Self.DNSName)

	return nil
}

// Listen returns a network listener on the Tailscale network
func (a *App) Listen(network, addr string) (net.Listener, error) {
	return a.server.Listen(network, addr)
}

// Close shuts down the Tailscale connection
func (a *App) Close() error {
	if a.server != nil {
		return a.server.Close()
	}
	return nil
}

// WhoIs performs a Tailscale WhoIs lookup for the given HTTP request
// and returns user information about the authenticated Tailscale user
func (a *App) WhoIs(r *http.Request) (*UserInfo, error) {
	if a.lc == nil {
		return nil, fmt.Errorf("local client not initialized")
	}

	// Extract IP address from RemoteAddr (format is typically "IP:Port")
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If no port, use RemoteAddr as-is
		remoteIP = r.RemoteAddr
	}

	fmt.Printf("[DEBUG] WhoIs lookup for: %s (original: %s)\n", remoteIP, r.RemoteAddr)

	info, err := a.lc.WhoIs(r.Context(), remoteIP)
	if err != nil {
		return nil, fmt.Errorf("whois lookup for %s: %w", remoteIP, err)
	}

	if info.Node == nil || info.UserProfile == nil {
		return nil, fmt.Errorf("incomplete whois response")
	}

	return &UserInfo{
		DisplayName:   info.UserProfile.DisplayName,
		LoginName:     info.UserProfile.LoginName,
		ProfilePicURL: info.UserProfile.ProfilePicURL,
		NodeName:      info.Node.ComputedName,
		NodeID:        string(info.Node.StableID),
	}, nil
}
