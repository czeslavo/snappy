package adapters

import (
	"errors"

	"github.com/czeslavo/snappy/internal/domain"
)

type CameraMock struct {
	Snapshots []domain.LoadedSnapshot
}

func (c *CameraMock) Get() (domain.LoadedSnapshot, error) {
	if len(c.Snapshots) < 1 {
		return domain.LoadedSnapshot{}, errors.New("no mocked snapshot")
	}

	var next domain.LoadedSnapshot
	next, c.Snapshots = c.Snapshots[len(c.Snapshots)-1], c.Snapshots[:len(c.Snapshots)-1]

	return next, nil
}
