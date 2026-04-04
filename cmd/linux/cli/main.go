package main

import (
	"os"

	"github.com/vasili-sikora/weather-app/internal/adapters/weather"
	"github.com/vasili-sikora/weather-app/internal/pkg/app/cli"
	"github.com/vasili-sikora/weather-app/internal/pkg/config"
	"github.com/vasili-sikora/weather-app/internal/pkg/flags"
	"github.com/vasili-sikora/weather-app/pkg/logger"
)

func main() {
	arguments := flags.Parse()

	configFile, err := os.Open(arguments.Path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := configFile.Close(); err != nil {
			panic(err)
		}
	}()

	appConfig, err := config.Parse(configFile)
	if err != nil {
		panic(err)
	}

	debugMode := os.Getenv("DEBUG") == "1"
	appLogger := logger.New(debugMode)
	weatherInfo := getProvider(appConfig, appLogger)

	app := cli.New(appLogger, weatherInfo, appConfig)
	if err := app.Run(); err != nil {
		appLogger.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func getProvider(appConfig config.Config, appLogger cli.Logger) cli.WeatherInfo {
	switch appConfig.P.Type {
	case "open-meteo":
		return weather.New(appLogger)
	default:
		return weather.New(appLogger)
	}
}
