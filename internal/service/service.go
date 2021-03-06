package service

import (
	"context"
	"sync"

	"github.com/czeslavo/snappy/internal/service/config"

	"github.com/sirupsen/logrus"

	"github.com/czeslavo/snappy/internal/ports"
)

type Service struct {
	HttpServer *ports.HTTPServer
	Ticker     *ports.Ticker
	Logger     logrus.FieldLogger
	Config     config.Config

	wg sync.WaitGroup
}

func (s *Service) Run(ctx context.Context) {
	s.Logger.WithField("config", s.Config).Info("Starting service...")

	s.runTicker(ctx)
	s.runHTTPServer(ctx)

	s.Logger.Info("Service running")
	s.wg.Wait()
}

func (s *Service) runTicker(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		if err := s.Ticker.Run(ctx); err != nil {
			s.Logger.WithError(err).Error("Ticker exited with error")
		} else {
			s.Logger.Info("Ticker exited")
		}
	}()
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		s.Ticker.Close()
	}()
}

func (s *Service) runHTTPServer(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil {
			s.Logger.WithError(err).Error("HTTP server exited with error")
		} else {
			s.Logger.Info("HTTP server exited")
		}
	}()
	go func() {
		defer s.wg.Done()
		<-ctx.Done()
		_ = s.HttpServer.Close()
	}()
}
