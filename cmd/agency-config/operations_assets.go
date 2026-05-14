package main

import (
	_ "embed"
	"net/http"
)

//go:embed operations_admin.js
var operationsAdminJS string

func (h *handler) operationsAsset(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin/operations/assets/operations.js" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write([]byte(operationsAdminJS))
}
