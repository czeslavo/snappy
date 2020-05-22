package adapters

import (
	"errors"

	"github.com/czeslavo/snappy/internal/domain"
)

type CameraMock struct {
	Snapshots []domain.Snapshot
}

func (c *CameraMock) Get() (domain.Snapshot, error) {
	if len(c.Snapshots) < 1 {
		return domain.Snapshot{}, errors.New("no mocked snapshot")
	}

	var next domain.Snapshot
	next, c.Snapshots = c.Snapshots[len(c.Snapshots)-1], c.Snapshots[:len(c.Snapshots)-1]

	return next, nil
}
