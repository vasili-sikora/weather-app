package cli

import (
	"fmt"

	"github.com/vasili-sikora/weather-app/internal/domain/models"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type WeatherInfo interface {
	GetTemperature(float64, float64) models.TempInfo
}

type cliApp struct {
	logger      Logger
	weatherInfo WeatherInfo
}

func New(logger Logger, weatherInfo WeatherInfo) *cliApp {
	return &cliApp{
		logger:      logger,
		weatherInfo: weatherInfo,
	}
}

func (c *cliApp) Run() error {
	const (
		latitude  = 53.6688
		longitude = 23.8223
	)

	tempInfo := c.weatherInfo.GetTemperature(latitude, longitude)
	c.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", tempInfo.Temp))

	return nil
}
