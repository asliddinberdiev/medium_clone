package service

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerWrite() io.Writer {
	_ = os.Mkdir("logs", 0770)

	logFilePath := path.Join("logs/", "access.log")

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Println("service_logger:  write - open file error: ", err)
	}
	return io.MultiWriter(logFile, os.Stdout)
}

func CustomLogger(out io.Writer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := ctx.ClientIP()

		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		statusCode := ctx.Writer.Status()

		formattedStart := start.Format("2006-01-02 15:04:05")

		logLine := fmt.Sprintf("%s | %s | %d | %v | %s  %s\n", formattedStart, clientIP, statusCode, latency, method, path)
		_, _ = out.Write([]byte(logLine))
	}
}
