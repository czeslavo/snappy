package application

import (
	"context"

	"github.com/czeslavo/snappy/internal/domain"
)

type LatestSnapshotRepo interface {
	GetLatest() (domain.LoadedSnapshot, error)
}

type GetLatestSnapshotHandler struct {
	repo LatestSnapshotRepo
}

func NewGetLatestSnapshotHandler(repo LatestSnapshotRepo) GetLatestSnapshotHandler {
	return GetLatestSnapshotHandler{repo}
}

func (h GetLatestSnapshotHandler) Execute(ctx context.Context) (domain.LoadedSnapshot, error) {
	return h.repo.GetLatest()
}
