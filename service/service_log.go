package service

import (
	"io"
	"log"
	"os"
	"path"
)

func LoggerWrite() io.Writer {
	_ = os.Mkdir("logs", 0770)

	logFilePath := path.Join("logs/", "access.log")

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	log.Println("service: logger open file error: ", err)
	return io.MultiWriter(logFile, os.Stdout)
}
