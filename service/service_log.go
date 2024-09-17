package service

import (
	"io"
	"os"
	"path"
)

func LoggerWrite() io.Writer {
	_ = os.Mkdir("logs", 0770)

	logFilePath := path.Join("logs/", "access.log")

	logFile, _ := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	return io.MultiWriter(logFile, os.Stdout)
}
