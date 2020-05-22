package application_test

import (
	"context"
	"image"
	"testing"
	"time"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/application"
	"github.com/czeslavo/snappy/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTakeSnapshotHandler(t *testing.T) {
	now := time.Now()
	expectedSnap := domain.MustNewSnapshot(
		image.NewGray(image.Rect(0, 0, 10, 10)),
		now,
	)
	camera := &adapters.CameraMock{
		Snapshots: []domain.Snapshot{expectedSnap},
	}
	repo := &adapters.SnapshotsInMemoryRepository{}
	handler := application.NewTakeSnapshotHandler(camera, repo)

	err := handler.Handle(context.Background(), now)
	require.NoError(t, err)

	snap, err := repo.GetLatest()
	require.NoError(t, err)
	require.Equal(t, expectedSnap, snap)
}

func TestTakeSnapshotHandler_cannot_get_snapshot(t *testing.T) {
	now := time.Now()
	camera := &adapters.CameraMock{}
	repo := &adapters.SnapshotsInMemoryRepository{}
	handler := application.NewTakeSnapshotHandler(camera, repo)

	err := handler.Handle(context.Background(), now)
	require.EqualError(t, err, "could not get snapshot: no mocked snapshot")

	_, err = repo.GetLatest()
	require.Error(t, err)
}
