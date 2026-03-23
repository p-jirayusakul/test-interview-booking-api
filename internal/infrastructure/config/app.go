package config

import "fmt"

type AppConfig struct {
	Port    string `mapstructure:"port"`
	Env     string `mapstructure:"env"`
	Host    string `mapstructure:"host"`
	BaseURL string `mapstructure:"base_url"`
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
}

func (s AppConfig) Validate() error {
	if s.Port == "" {
		return fmt.Errorf("app.port is required")
	}
	if s.Env == "" {
		return fmt.Errorf("app.env is required")
	}
	if s.Host == "" {
		return fmt.Errorf("app.host is required")
	}
	if s.BaseURL == "" {
		return fmt.Errorf("app.base_url is required")
	}
	if s.Version == "" {
		return fmt.Errorf("app.version is required")
	}
	if s.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	return nil
}
