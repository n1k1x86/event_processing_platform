package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	App         App         `json:"app"`
	HTTPServer  HTTPServer  `json:"http_server"`
	PprofServer PprofServer `json:"pprof_server"`
	Jobs        Jobs        `json:"jobs"`
}

type App struct {
	Debug           bool     `json:"debug"`
	GracefulTimeout Duration `json:"graceful_timeout"`
}

type HTTPServer struct {
	Addr         string   `json:"addr"`
	WriteTimeout Duration `json:"write_timeout"`
	ReadTimeout  Duration `json:"read_timeout"`
}

type PprofServer struct {
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

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		dur, err := time.ParseDuration(s)
		fmt.Println(dur)
		if err != nil {
			return err
		}
		d.Duration = dur
		return nil
	}

	var ns int64
	if err := json.Unmarshal(b, &ns); err != nil {
		return err
	}
	d.Duration = time.Duration(ns)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
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
