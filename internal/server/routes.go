package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	mux.HandleFunc("/video/{file}", s.VideoHandler)

	mux.HandleFunc("/health", s.healthHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) VideoHandler(w http.ResponseWriter, r *http.Request) {
	file := r.PathValue("file")
	fmt.Println(file)
	switch path.Ext(file) {
	case ".ts":
		w.Header().Add("Content-Type", "video/MP2T")
		http.ServeFile(w, r, path.Join(".", "cache", "vid", file))
	case ".m3u8":
		w.Header().Add("Content-Type", "application/vnd.apple.mpegurl")
		http.ServeFile(w, r, path.Join(".", "cache", "vid", file))
	default:
		fmt.Println(path.Ext(file))
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
