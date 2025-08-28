package server

import (
	"errors"
	"fmt"
	"hls-on-the-fly/internal/ffmpeg"
	"hls-on-the-fly/internal/m3u8"
	pathhelpers "hls-on-the-fly/internal/path_helpers"
	"net/http"
	"os"
	"path"
	"sync"
)

var mu sync.Mutex

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	mux.HandleFunc("/video/{file}", s.VideoHandler)

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
	p := path.Join(".", "cache", "vid", file)
	fileWithoutExtension := pathhelpers.GetNameWithoutExtension(file)

	switch path.Ext(file) {
	case ".ts":
		mu.Lock()
		if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
			baseFileName := pathhelpers.GetNameWithoutExtension(fileWithoutExtension)
			manifestFilePath := path.Join(".", "cache", "vid", fmt.Sprintf("%v.m3u8", baseFileName))

			segments, err := m3u8.ParseManifest(manifestFilePath)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("could not parse manifest file"))
				return
			}

			var segment m3u8.Segment
			for _, s := range segments {
				if s.Name == file {
					segment = s
				}
			}

			if s == nil {
				w.WriteHeader(500)
				w.Write([]byte("no segment exists for request"))
				return
			}

			_, err = ffmpeg.HLSChunk(int(segment.Duration), int(segment.Start), path.Join(".", "tmp", baseFileName+".mp4"), p)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("could not transcode segment"))
				return
			}
		}

		w.Header().Add("Content-Type", "video/MP2T")
		http.ServeFile(w, r, path.Join(".", "cache", "vid", file))
		mu.Unlock()
	case ".m3u8":
		if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
			videoFile := path.Join(".", "tmp", fmt.Sprintf("%v.mp4", fileWithoutExtension))
			out, err := m3u8.CreateManifestForFile(videoFile, s.env.HlsTime, s.env.Cache)
			if err != nil {
				fmt.Println("could not generate manifest file: ", err.Error())

				w.WriteHeader(500)
				w.Write([]byte("could not generate manifest"))
			}

			fmt.Println("generated manifest: ", out)
		}

		w.Header().Add("Content-Type", "application/vnd.apple.mpegurl")
		http.ServeFile(w, r, p)
	default:
		fmt.Println(path.Ext(file))
	}
}
