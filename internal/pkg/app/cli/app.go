package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type cliApp struct {
	logger Logger
}

func New(logger Logger) *cliApp {
	return &cliApp{logger: logger}
}

func (c *cliApp) Run() error {
	type Current struct {
		Temp float32 `json:"temperature_2m"`
	}

	type Response struct {
		Curr Current `json:"current"`
	}

	var response Response
	params := fmt.Sprintf("latitude=%f&longitude=%f&current=temperature_2m", 53.6688, 23.8223)
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)

	c.logger.Debug(fmt.Sprintf("request url: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		customErr := errors.New("can't get weather data from openmeteo")
		c.logger.Error(customErr.Error())
		return errors.Join(customErr, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error(fmt.Sprintf("can't close body err - %s", err.Error()))
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		customErr := errors.New("can't read data from response")
		c.logger.Error(customErr.Error())
		return errors.Join(customErr, err)
	}

	if err := json.Unmarshal(data, &response); err != nil {
		customErr := errors.New("can't unmarshal data from response")
		c.logger.Error(customErr.Error())
		return errors.Join(customErr, err)
	}

	c.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp))
	return nil
}
