package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Loader struct {
	v *viper.Viper
}

func NewLoader() *Loader {
	v := viper.New()
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &Loader{v: v}
}

func (l *Loader) Load(path string) error {
	l.v.SetConfigFile(path)
	return l.v.ReadInConfig()
}

func (l *Loader) App() (AppConfig, error) {
	cfg := AppConfig{
		Port:    l.v.GetString("app.port"),
		Env:     l.v.GetString("app.env"),
		Host:    l.v.GetString("app.host"),
		BaseURL: l.v.GetString("app.base_url"),
		Version: l.v.GetString("app.version"),
		Name:    l.v.GetString("app.name"),
	}

	err := cfg.Validate()
	if err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}

func (l *Loader) Database() (DatabaseConfig, error) {
	cfg := DatabaseConfig{
		Host:     l.v.GetString("db.host"),
		Port:     l.v.GetString("db.port"),
		User:     l.v.GetString("db.user"),
		Password: l.v.GetString("db.password"),
		DBName:   l.v.GetString("db.dbname"),
	}

	return cfg, nil
}
