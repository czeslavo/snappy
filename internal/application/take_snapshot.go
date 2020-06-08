package application

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/czeslavo/snappy/internal/domain"
)

type Camera interface {
	Get() (domain.LoadedSnapshot, error)
}

type SnapshotsRepository interface {
	Store(snapshot domain.LoadedSnapshot) error
}

type TakeSnapshotHandler struct {
	camera        Camera
	snapshotsRepo SnapshotsRepository
	logger        logrus.FieldLogger
}

func NewTakeSnapshotHandler(camera Camera, repository SnapshotsRepository, logger logrus.FieldLogger) TakeSnapshotHandler {
	return TakeSnapshotHandler{
		camera,
		repository,
		logger,
	}
}

func (h TakeSnapshotHandler) Handle(_ context.Context, _ time.Time) error {
	h.logger.Info("Taking snapshot...")

	snapshot, err := h.camera.Get()
	if err != nil {
		return fmt.Errorf("could not get snapshot: %s", err)
	}

	if err := h.snapshotsRepo.Store(snapshot); err != nil {
		return fmt.Errorf("could not store snapshot: %s", err)
	}

	return nil
}
