// Package logger provides functions to initialize and access the application logger.
package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger

// InitLogger initializes the application logger.
func InitLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	Log = logger.Sugar()
	return nil
}
