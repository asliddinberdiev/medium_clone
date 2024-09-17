package service

import (
	L "log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(dirName, fileName string) *zap.Logger {

	logFilePath := filepath.Join(dirName, fileName+".log")

	// Create logs directory if it doesn't exist
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			L.Fatalf("failed to create log directory: %v", err)
		}
	}

	// Create app.log if they don't exist
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			L.Fatalf("failed to create log file: %v", err)
		}
		file.Close()
	}

	// Configure the logger
	config := zap.Config{
		Encoding:         zap.NewProductionConfig().Encoding,
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{logFilePath, "stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			LevelKey:    "level",
			MessageKey:  "msg",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		L.Panicf("create logger error: %v\n", err)
	}

	return logger
}
