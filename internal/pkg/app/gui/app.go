package gui

import (
	"fmt"

	guisettings "github.com/vasili-sikora/weather-app/internal/domain/gui_settings"
	"github.com/vasili-sikora/weather-app/internal/domain/models"
	"github.com/vasili-sikora/weather-app/internal/pkg/config"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string, err error)
}

type WeatherInfo interface {
	GetTemperature(float64, float64) (models.TempInfo, error)
}

type Provider interface {
	CreateWindow(name string, size guisettings.WindowSize) (guisettings.Window, error)
	GetAppRunner() guisettings.AppRunner
	GetTextWidget(text string) guisettings.TextWidget
}

type guiApp struct {
	logger      Logger
	provider    Provider
	weatherInfo WeatherInfo
	config      config.Config
}

func New(logger Logger, provider Provider, weatherInfo WeatherInfo, appConfig config.Config) *guiApp {
	return &guiApp{
		logger:      logger,
		provider:    provider,
		weatherInfo: weatherInfo,
		config:      appConfig,
	}
}

func (g *guiApp) Run() error {
	window, err := g.provider.CreateWindow("Weather informer", guisettings.NewWS(400, 200))
	if err != nil {
		return err
	}

	label := g.provider.GetTextWidget("Загрузка температуры...")
	if err := window.SetTemperatureWidget(label); err != nil {
		return err
	}

	tempInfo, err := g.weatherInfo.GetTemperature(g.config.L.Lat, g.config.L.Long)
	if err != nil {
		g.logger.Error("can't get temp info", err)
		return err
	}

	if err := window.UpdateTemperature(tempInfo.Temp); err != nil {
		return err
	}

	if err := window.Render(); err != nil {
		return err
	}

	g.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", tempInfo.Temp))
	g.provider.GetAppRunner().Run()

	return nil
}
