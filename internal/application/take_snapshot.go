package application

import (
	"context"
	"fmt"
	"time"

	"github.com/czeslavo/snappy/internal/domain"
)

type Camera interface {
	Get() (domain.Snapshot, error)
}

type SnapshotsRepository interface {
	Store(snapshot domain.Snapshot) error
}

type TakeSnapshotHandler struct {
	camera        Camera
	snapshotsRepo SnapshotsRepository
}

func NewTakeSnapshotHandler(camera Camera, repository SnapshotsRepository) TakeSnapshotHandler {
	return TakeSnapshotHandler{
		camera:        camera,
		snapshotsRepo: repository,
	}
}

func (h TakeSnapshotHandler) Handle(_ context.Context, _ time.Time) error {
	snapshot, err := h.camera.Get()
	if err != nil {
		return fmt.Errorf("could not get snapshot: %s", err)
	}

	if err := h.snapshotsRepo.Store(snapshot); err != nil {
		return fmt.Errorf("could not store snapshot: %s", err)
	}

	return nil
}
