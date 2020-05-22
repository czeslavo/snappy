package ports

import (
	"image/jpeg"
	"net/http"

	"github.com/czeslavo/snappy/internal/application"
)

type HTTPServer struct {
	mux *http.ServeMux
	srv *http.Server

	getLatestSnapshotHandler application.GetLatestSnapshotHandler
}

func NewHTTPServer(getLatestSnapshotHandler application.GetLatestSnapshotHandler) *HTTPServer {
	s := &HTTPServer{
		mux:                      http.NewServeMux(),
		getLatestSnapshotHandler: getLatestSnapshotHandler,
	}

	s.mux.Handle("/latest.jpeg", http.HandlerFunc(s.handleLatest))

	s.srv = &http.Server{Handler: s.mux}

	return s
}

func (s *HTTPServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Close() error {
	return s.srv.Close()
}

func (s *HTTPServer) handleLatest(resp http.ResponseWriter, req *http.Request) {
	snapshot, err := s.getLatestSnapshotHandler.Execute(req.Context())
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := jpeg.Encode(resp, snapshot.Image(), nil); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}
