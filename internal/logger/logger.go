package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// What does this file do: returns a custom zap instance with human readable time
// (ISO format) instead of UNIX timestamp(epoch)

// To change the timestamp format from epoch to ISO format
func GetLogger() *zap.Logger{
	// Step 1: Start with the default/base config
	config := zap.NewProductionConfig()

	// todo: make GetLogger() handle the prod and dev environment by using .env

	// Step 2: Do the necessary changes
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Step 3: Build the logger 
	logger, err := config.Build() // Build constructs a logger from the Config and Options
	if err != nil { // i.e error illama illa -> error irukku
		panic(err) //kinda similar to exit(1), but different 
		// panic() is a recoverable crash that runs cleanup tasks (defer),
		// while os.Exit() is an immediate shutdown that stops everything instantly.
	}
	return logger
}
