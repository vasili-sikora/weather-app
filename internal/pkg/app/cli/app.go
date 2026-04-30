package cli

import (
	"fmt"

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

type cliApp struct {
	logger      Logger
	weatherInfo WeatherInfo
	config      config.Config
}

func New(logger Logger, weatherInfo WeatherInfo, appConfig config.Config) *cliApp {
	return &cliApp{
		logger:      logger,
		weatherInfo: weatherInfo,
		config:      appConfig,
	}
}

func (c *cliApp) Run() error {
	tempInfo, err := c.weatherInfo.GetTemperature(c.config.L.Lat, c.config.L.Long)
	if err != nil {
		c.logger.Error("can't get temp info", err)
		return err
	}

	c.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", tempInfo.Temp))

	return nil
}
