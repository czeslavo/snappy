// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package service

import (
	"github.com/czeslavo/snappy/internal/adapters"
	"github.com/czeslavo/snappy/internal/application"
	"github.com/czeslavo/snappy/internal/ports"
	"github.com/czeslavo/snappy/internal/service/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Injectors from wire.go:

func BuildService() (*Service, error) {
	configConfig, err := config.ReadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	httpPort := configConfig.HTTPPort
	snapshotsDirectory := configConfig.SnapshotsDirectory
	snapshotsFileSystemRepository, err := adapters.NewSnapshotsFileSystemRepository(snapshotsDirectory)
	if err != nil {
		return nil, err
	}
	getLatestSnapshotHandler := application.NewGetLatestSnapshotHandler(snapshotsFileSystemRepository)
	fieldLogger := provideLogger()
	httpServer := ports.NewHTTPServer(httpPort, getLatestSnapshotHandler, fieldLogger)
	client := _wireClientValue
	cameraURL := configConfig.CameraURL
	jpegCamera, err := adapters.NewJPEGCamera(client, cameraURL)
	if err != nil {
		return nil, err
	}
	takeSnapshotHandler := application.NewTakeSnapshotHandler(jpegCamera, snapshotsFileSystemRepository, fieldLogger)
	zipSnapshotsArchiver := adapters.NewZipSnapshotArchiver()
	ftpUploader := provideFtpUploader(configConfig, fieldLogger)
	archiveAllSnapshotsHandler := application.NewArchiveAllSnapshotsHandler(snapshotsFileSystemRepository, zipSnapshotsArchiver, ftpUploader, fieldLogger)
	ticker := ports.NewTicker(takeSnapshotHandler, archiveAllSnapshotsHandler, configConfig, fieldLogger)
	service := &Service{
		HttpServer: httpServer,
		Ticker:     ticker,
		Logger:     fieldLogger,
		Config:     configConfig,
	}
	return service, nil
}

var (
	_wireClientValue = &http.Client{}
)

// wire.go:

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

func provideFtpUploader(conf config.Config, logger logrus.FieldLogger) adapters.FtpUploader {
	return adapters.NewFtpUploader(adapters.Credentials{
		Username: conf.FtpUsername,
		Password: conf.FtpPassword,
	}, conf.FtpHost, conf.FtpTargetDirectory, logger)
}
