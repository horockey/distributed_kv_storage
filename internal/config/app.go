package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed app.json
var cfgData []byte

const AppName string = "distributed_kv_storage"

type Config struct {
	InstanceManager InstanceManager `json:"instance_manager"`
	Http            Http            `json:"http"`
	Hostname        string          `json:"-"`
}

type InstanceManager struct {
	PollIntervalMsec     int `json:"poll_interval_msec"`
	DownHoldDirationMsec int `json:"down_hold_duration_msec"`
}

type Http struct {
	Port     int    `json:"port"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshalling cfg data json: %w", err)
	}

	var err error
	cfg.Hostname, err = os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostnane: %w", err)
	}

	return &cfg, nil
}
