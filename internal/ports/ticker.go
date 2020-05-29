package ports

import (
	"context"
	"sync"
	"time"

	"github.com/czeslavo/snappy/internal/application"
	"github.com/czeslavo/snappy/internal/service/config"

	"github.com/sirupsen/logrus"
)

type TickerHandler struct {
	Frequency time.Duration
	Handler   TickerHandlerFunc
}

type TickerHandlerFunc func(ctx context.Context, now time.Time) error

type Ticker struct {
	takeSnapshotHandler TickerHandler

	logger logrus.FieldLogger
	close  chan struct{}
	wg     sync.WaitGroup
}

func NewTicker(takeSnaphotHandler application.TakeSnapshotHandler, conf config.Config, logger logrus.FieldLogger) *Ticker {
	return &Ticker{
		takeSnapshotHandler: TickerHandler{
			Frequency: time.Duration(conf.SnapshotsFrequency),
			Handler:   takeSnaphotHandler.Handle,
		},
		logger: logger,
		close:  make(chan struct{}),
	}
}

func (t *Ticker) Run(ctx context.Context) error {
	t.handleWithTicker(ctx, t.takeSnapshotHandler.Frequency, t.takeSnapshotHandler.Handler)

	t.wg.Wait()
	return nil
}

func (t *Ticker) handleWithTicker(ctx context.Context, tick time.Duration, handler TickerHandlerFunc) {
	tickFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				t.logger.Errorf("Tick panicked: %v", r)
			}
		}()

		if err := handler(ctx, time.Now()); err != nil {
			t.logger.WithError(err).Error("Tick failed")
		}
	}

	t.wg.Add(1)
	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				tickFunc()
			case <-ctx.Done():
				t.wg.Done()
				return
			case <-t.close:
				t.wg.Done()
				return
			}
		}
	}()
}

func (t *Ticker) Close() {
	close(t.close)
}
