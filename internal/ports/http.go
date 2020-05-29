package ports

import (
	"fmt"
	"image/jpeg"
	"net/http"

	"github.com/czeslavo/snappy/internal/service/config"

	"github.com/sirupsen/logrus"

	"github.com/czeslavo/snappy/internal/application"
)

type HTTPServer struct {
	mux    *http.ServeMux
	srv    *http.Server
	logger logrus.FieldLogger

	getLatestSnapshotHandler application.GetLatestSnapshotHandler
}

func NewHTTPServer(port config.HTTPPort, getLatestSnapshotHandler application.GetLatestSnapshotHandler, logger logrus.FieldLogger) *HTTPServer {
	s := &HTTPServer{
		mux:                      http.NewServeMux(),
		getLatestSnapshotHandler: getLatestSnapshotHandler,
		logger:                   logger,
	}

	s.mux.Handle("/snapshots/latest.jpeg", http.HandlerFunc(s.handleLatest))

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.mux,
	}

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

	resp.Header().Add("taken_at", snapshot.TakenAt().String())
	if err := jpeg.Encode(resp, snapshot.Image(), nil); err != nil {
		s.logger.WithError(err).Warning("Snapshot encoding failed")
		return
	}
}
