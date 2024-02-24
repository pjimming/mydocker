package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func Log(containerId string) error {
	logFilePath := filepath.Join(getContainerDir(containerId), LogFile)
	file, err := os.Open(logFilePath)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		logrus.Errorf("[Log] open %s error, %v", logFilePath, err)
		return err
	}
	content, err := io.ReadAll(file)
	if err != nil {
		logrus.Errorf("[Log] read file %s error %v", logFilePath, err)
		return err
	}
	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		logrus.Errorf("[Log] Fprint error %v", err)
		return err
	}
	return nil
}
