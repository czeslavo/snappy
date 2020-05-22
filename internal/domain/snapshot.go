package domain

import (
	"errors"
	"image"
	"time"
)

type Snapshot struct {
	img     image.Image
	takenAt time.Time
}

func NewSnapshot(img image.Image, takenAt time.Time) (Snapshot, error) {
	if img == nil {
		return Snapshot{}, errors.New("empty image")
	}

	return Snapshot{
		img:     img,
		takenAt: takenAt,
	}, nil
}

func MustNewSnapshot(img image.Image, takenAt time.Time) Snapshot {
	s, err := NewSnapshot(img, takenAt)
	if err != nil {
		panic(err)
	}
	return s
}

func (s Snapshot) Image() image.Image {
	return s.img
}

func (s Snapshot) TakenAt() time.Time {
	return s.takenAt
}
