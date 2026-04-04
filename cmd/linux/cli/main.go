package main

import (
	"os"

	"github.com/vasili-sikora/weather-app/internal/adapters/weather"
	"github.com/vasili-sikora/weather-app/internal/pkg/app/cli"
	"github.com/vasili-sikora/weather-app/pkg/logger"
)

func main() {
	debugMode := os.Getenv("DEBUG") == "1"
	appLogger := logger.New(debugMode)
	weatherInfo := weather.New(appLogger)

	app := cli.New(appLogger, weatherInfo)
	if err := app.Run(); err != nil {
		appLogger.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
