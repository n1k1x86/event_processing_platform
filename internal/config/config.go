package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App        App        `json:"app"`
	HTTPServer HTTPServer `json:"http_server"`
}

type App struct {
	Debug bool `json:"debug"`
}
type HTTPServer struct {
	Addr string `json:"addr"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config/dev.json")
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
