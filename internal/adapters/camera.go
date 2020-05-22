package adapters

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/czeslavo/snappy/internal/domain"
)

type JPEGCamera struct {
	client *http.Client
	url    string
}

func NewJPEGCamera(client *http.Client, url string) (JPEGCamera, error) {
	if client == nil {
		return JPEGCamera{}, errors.New("empty client")
	}
	if url == "" {
		return JPEGCamera{}, errors.New("empty url")
	}

	return JPEGCamera{
		client: client,
		url:    url,
	}, nil
}

func (c JPEGCamera) Get() (domain.Snapshot, error) {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("could not get snapshot: %s", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return domain.Snapshot{}, fmt.Errorf("status code %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return domain.Snapshot{}, fmt.Errorf("could not decode image: %s", err)
	}

	return domain.NewSnapshot(img, time.Now())
}
