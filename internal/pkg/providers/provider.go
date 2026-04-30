package providers

import (
	pogodaby "github.com/vasili-sikora/weather-app/internal/adapters/pogoda_by"
	"github.com/vasili-sikora/weather-app/internal/adapters/weather"
	"github.com/vasili-sikora/weather-app/internal/pkg/app/cli"
	"github.com/vasili-sikora/weather-app/internal/pkg/config"
)

func GetProvider(c config.Config, l cli.Logger) cli.WeatherInfo {
	switch c.P.Type {
	case "open-meteo":
		return weather.New(l)
	case "pogoda":
		return pogodaby.New(l)
	default:
		return weather.New(l)
	}
}
