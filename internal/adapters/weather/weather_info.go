package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vasili-sikora/weather-app/internal/domain/models"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

type current struct {
	Temp float32 `json:"temperature_2m"`
}

type response struct {
	Curr current `json:"current"`
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string, err error)
}

type weatherInfo struct {
	current  current
	logger   Logger
	isLoaded bool
}

func New(logger Logger) *weatherInfo {
	return &weatherInfo{
		logger: logger,
	}
}

func (wi *weatherInfo) getWeatherInfo(lat, long float64) error {
	if wi.isLoaded {
		return nil
	}

	var weatherResponse response

	params := fmt.Sprintf(
		"latitude=%f&longitude=%f&current=temperature_2m",
		lat,
		long,
	)
	url := fmt.Sprintf("%s?%s", apiURL, params)

	wi.logger.Debug(fmt.Sprintf("request url: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		wi.logger.Error("can't get weather data", err)
		customErr := errors.New("can't get weather data from openmeteo")
		return errors.Join(customErr, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			wi.logger.Error("can't close body", err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		wi.logger.Error("can't read data from body", err)
		customErr := errors.New("can't read data from response")
		return errors.Join(customErr, err)
	}

	if err := json.Unmarshal(data, &weatherResponse); err != nil {
		wi.logger.Error("can't unmarshal json data", err)
		customErr := errors.New("can't unmarshal data from response")
		return errors.Join(customErr, err)
	}

	wi.current = weatherResponse.Curr
	wi.isLoaded = true

	return nil
}

func (wi *weatherInfo) GetTemperature(lat, long float64) (models.TempInfo, error) {
	err := wi.getWeatherInfo(lat, long)

	return models.TempInfo{
		Temp: wi.current.Temp,
	}, err
}
