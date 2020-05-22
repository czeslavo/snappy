package ports

import (
	"context"
	"sync"
	"time"

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
	wg     sync.WaitGroup
}

func NewTicker(takeSnapshotHandler TickerHandler, logger logrus.FieldLogger) *Ticker {
	return &Ticker{
		takeSnapshotHandler: takeSnapshotHandler,
		logger:              logger,
	}
}

func (t *Ticker) Run(ctx context.Context) error {
	t.handleWithTicker(ctx, t.takeSnapshotHandler.Frequency, t.takeSnapshotHandler.Handler)

	t.wg.Wait()
	return nil
}

func (t *Ticker) handleWithTicker(ctx context.Context, tick time.Duration, handler TickerHandlerFunc) {
	t.wg.Add(1)
	go func() {
		for {
			select {
			case <-time.Tick(tick):
				if err := handler(ctx, time.Now()); err != nil {
					t.logger.WithError(err).Error("Tick failed")
				}
			case <-ctx.Done():
				t.wg.Done()
				return
			}
		}
	}()
}
