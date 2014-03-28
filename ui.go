package main

import (
	"net/http"
)

// homeRedirectHandler exists simply to redirect root requests to /ui
func homeRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/index.html", http.StatusMovedPermanently)
}
