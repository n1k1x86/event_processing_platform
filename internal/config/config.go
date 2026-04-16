package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App        App        `json:"app"`
	HTTPServer HTTPServer `json:"http_server"`
	Jobs       Jobs       `json:"jobs"`
}

type App struct {
	Debug bool `json:"debug"`
}

type HTTPServer struct {
	Addr string `json:"addr"`
}

type Jobs struct {
	SendEmailJob Job `json:"send_email_job"`
}

type Job struct {
	Workers   int    `json:"workers"`
	QueueSize int    `json:"queue_size"`
	JobType   string `json:"job_type"`
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
