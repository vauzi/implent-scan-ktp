package web

import (
	"net/http"

	s "github.com/vauzi/implent-scan-ktp/server"
)

func Route() http.Handler {
	s := &s.Server{}
	mux := http.NewServeMux()

	// api routes get
	mux.HandleFunc("/api/health-check", s.HealthCheck)

	// api routes post
	mux.HandleFunc("/api/file-upload", s.UploadHandler)

	mux.HandleFunc("/api/parse-nik", s.ParseNiks)

	return mux
}
