package web

import "embed"

// FrontendFS contains the embedded React frontend build output.
// The frontend is built by Vite into internal/web/frontend/dist/
//
//go:embed all:frontend/dist
var FrontendFS embed.FS
