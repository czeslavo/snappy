// +build wireinject

package service

import (
	"net/http"
	"os"

	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/application"
	"github.com/czeslavo/snappy/internal/ports"
	"github.com/czeslavo/snappy/internal/service/config"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

func BuildService() (*Service, error) {
	wire.Build(
		wire.Struct(new(Service),
			"HTTPServer",
			"Ticker",
			"Logger",
			"Config",
		),

		provideLogger,
		config.ReadConfigFromEnv,
		wire.FieldsOf(new(config.Config),
			"CameraURL",
			"SnapshotsDirectory",
			"HTTPPort",
		),

		ports.NewHTTPServer,
		ports.NewTicker,

		wire.Bind(new(application.LatestSnapshotRepo), new(adapters.SnapshotsFileSystemRepository)),
		wire.Bind(new(application.SnapshotsRepository), new(adapters.SnapshotsFileSystemRepository)),
		wire.Bind(new(application.Camera), new(adapters.JPEGCamera)),

		adapters.NewJPEGCamera,
		wire.Value(&http.Client{}),
		adapters.NewSnapshotsFileSystemRepository,

		application.NewGetLatestSnapshotHandler,
		application.NewTakeSnapshotHandler,
	)
	return &Service{}, nil
}

func provideLogger() logrus.FieldLogger {
	level := logrus.DebugLevel
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err == nil {
			level = l
		}
	}

	logger := logrus.New()
	logger.SetLevel(level)

	return logger
}
