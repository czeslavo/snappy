// +build integration

package adapters_test

import (
	"image"
	"net/http"
	"testing"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJPEGCamera(t *testing.T) {
	client := &http.Client{}
	c, err := adapters.NewJPEGCamera(client, "http://localhost:8085/sample.jpeg")
	require.NoError(t, err)

	snap, err := c.Get()
	require.NoError(t, err)
	require.NotNil(t, snap)

	assert.Equal(t, image.Point{X: 0, Y: 0}, snap.Image().Bounds().Min)
	assert.Equal(t, image.Point{X: 1920, Y: 1080}, snap.Image().Bounds().Max)
	assert.NotEmpty(t, snap.TakenAt())
}
