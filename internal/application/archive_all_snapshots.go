package application

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/czeslavo/snappy/internal/domain"
)

type AllSnapshotsRepository interface {
	GetAll() ([]domain.FileSnapshot, error)
}

type ArchivableSnapshot interface {
	TakenAt() time.Time
	Copy(io.Writer) error
}

type SnapshotsArchiver interface {
	Archive([]ArchivableSnapshot) (path string, err error)
}

type ArchiveUploader interface {
	Upload(path string) error
}

type ArchiveAllSnapshotsHandler struct {
	repo     AllSnapshotsRepository
	archiver SnapshotsArchiver
	uploader ArchiveUploader
	logger   logrus.FieldLogger
}

func NewArchiveAllSnapshotsHandler(repo AllSnapshotsRepository, archiver SnapshotsArchiver, uploader ArchiveUploader, logger logrus.FieldLogger) ArchiveAllSnapshotsHandler {
	return ArchiveAllSnapshotsHandler{
		repo,
		archiver,
		uploader,
		logger.WithField("handler", "archive_all_snapshots"),
	}
}

func (h ArchiveAllSnapshotsHandler) Handle(ctx context.Context, _ time.Time) error {
	h.logger.Info("Archiving...")

	snaps, err := h.repo.GetAll()
	if err != nil {
		return fmt.Errorf("could not get all snapshots: %s", err)
	}
	path, err := h.archiver.Archive(asArchivableSlice(snaps))
	if err != nil {
		return fmt.Errorf("could not archive all snapshots: %s", err)
	}
	defer os.Remove(path)
	h.logger.Debugf("Archived to %s", path)

	if err := h.uploader.Upload(path); err != nil {
		return fmt.Errorf("could not upload archive: %s", err)
	}
	h.logger.Debugf("Uploaded archive: %s", path)

	if err := removeAll(snaps); err != nil {
		return fmt.Errorf("could not remove snapshots")
	}

	return nil
}

func asArchivableSlice(snaps []domain.FileSnapshot) []ArchivableSnapshot {
	var archivable []ArchivableSnapshot
	for _, s := range snaps {
		archivable = append(archivable, s)
	}
	return archivable
}

func removeAll(snaps []domain.FileSnapshot) error {
	for _, s := range snaps {
		if err := os.Remove(s.Path()); err != nil {
			return fmt.Errorf("could not remove snapshot: %s", err)
		}
	}

	return nil
}
