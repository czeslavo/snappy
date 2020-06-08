package adapters

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jlaffaye/ftp"
)

type FtpUploader struct {
	creds     Credentials
	host      string
	targetDir string
	logger    logrus.FieldLogger
}

type Credentials struct {
	Username, Password string
}

func NewFtpUploader(creds Credentials, host, targetDir string, logger logrus.FieldLogger) FtpUploader {
	return FtpUploader{creds, host, targetDir, logger}
}

func (u FtpUploader) Upload(path string) error {
	conn, err := ftp.Dial(u.host, ftp.DialWithTimeout(time.Second*20))
	if err != nil {
		return fmt.Errorf("could not dial ftp: %s", err)
	}
	defer conn.Quit()

	if err := conn.Login(u.creds.Username, u.creds.Password); err != nil {
		return fmt.Errorf("could not login: %s", err)
	}
	defer conn.Logout()

	if err := conn.ChangeDir(u.targetDir); err != nil {
		return fmt.Errorf("could not change directory to target: %s", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %s", err)
	}
	defer f.Close()

	u.logger.Debug("Transferring file to FTP: %s", path)
	if err := conn.Stor(filepath.Base(path), f); err != nil {
		return fmt.Errorf("could not store file: %s", err)
	}

	return nil
}
