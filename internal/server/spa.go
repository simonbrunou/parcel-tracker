package server

import (
	"io/fs"
	"net/http"
	"strings"
)

// SPAHandler serves the embedded Svelte SPA.
// Real files (JS, CSS, images) are served directly with cache headers.
// All other paths get index.html (SPA client-side routing fallback).
func SPAHandler(distFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if the file exists
		if f, err := distFS.Open(path); err == nil {
			f.Close()

			// Serve hashed assets with aggressive caching
			if strings.HasPrefix(path, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}

			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html
		r.URL.Path = "/"
		w.Header().Set("Cache-Control", "no-cache")
		fileServer.ServeHTTP(w, r)
	})
}
