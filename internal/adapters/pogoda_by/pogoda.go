package pogodaby

import (
	"encoding/json"
	"net/http"

	"github.com/vasili-sikora/weather-app/internal/domain/models"
)

const url = "https://pogoda.by/api/v2/weather-fact?station=26820"

type resp struct {
	Temp float32 `json:"t"`
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string, err error)
}

type pogoda struct {
	logger Logger
}

func New(logger Logger) *pogoda {
	return &pogoda{logger: logger}
}

func (p *pogoda) GetTemperature(lat, long float64) (models.TempInfo, error) {
	response, err := http.Get(url)
	if err != nil {
		p.logger.Error("can't get data from pogoda.by", err)
		return models.TempInfo{}, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			p.logger.Error("can't close response body", err)
		}
	}()

	var weatherResponse resp
	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		p.logger.Error("can't decode JSON", err)
		return models.TempInfo{}, err
	}

	return models.TempInfo{
		Temp: weatherResponse.Temp,
	}, nil
}
