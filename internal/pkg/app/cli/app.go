package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type cliApp struct {
	logger Logger
	cache  Cache
}

func New(logger Logger) *cliApp {
	return &cliApp{logger: logger,
		cache: newMemoryCache()}
}

func (c *cliApp) Run() error {
	type Current struct {
		Temp float32 `json:"temperature_2m"`
	}

	type Response struct {
		Curr Current `json:"current"`
	}

	const (
		latitude  = 53.6688
		longitude = 23.8223
	)

	ctx := context.Background()
	cacheKey := fmt.Sprintf("weather:%f:%f", latitude, longitude)

	if cachedTemp, ok, err := c.cache.Get(ctx, cacheKey); err != nil {
		c.logger.Error(fmt.Sprintf("cache read failed: %s", err.Error()))
	} else if ok {
		c.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия (cache)", cachedTemp))
		return nil
	}
	var response Response
	params := fmt.Sprintf("latitude=%f&longitude=%f&current=temperature_2m", latitude, longitude)
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

	ttlSeconds := 60
	if ttlRaw := os.Getenv("CACHE_TTL_SECONDS"); ttlRaw != "" {
		if parsed, parseErr := strconv.Atoi(ttlRaw); parseErr == nil && parsed > 0 {
			ttlSeconds = parsed
		}
	}

	if err := c.cache.Set(ctx, cacheKey, response.Curr.Temp, time.Duration(ttlSeconds)*time.Second); err != nil {
		c.logger.Error(fmt.Sprintf("cache write failed: %s", err.Error()))
	}

	c.logger.Info(fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp))
	return nil
}
