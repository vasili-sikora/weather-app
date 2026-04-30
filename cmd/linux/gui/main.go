package main

import (
	"os"

	"github.com/vasili-sikora/weather-app/internal/pkg/app/gui"
	"github.com/vasili-sikora/weather-app/internal/pkg/config"
	"github.com/vasili-sikora/weather-app/internal/pkg/flags"
	fynegui "github.com/vasili-sikora/weather-app/internal/pkg/gui/fyne"
	"github.com/vasili-sikora/weather-app/internal/pkg/providers"
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
	weatherInfo := providers.GetProvider(appConfig, appLogger)
	provider := fynegui.NewP()
	app := gui.New(appLogger, provider, weatherInfo, appConfig)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
