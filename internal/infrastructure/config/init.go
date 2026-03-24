package config

type Config struct {
	AppCfg      AppConfig
	DatabaseCfg DatabaseConfig
	Loader      *Loader
}

func InitConfig() (*Config, error) {
	loader := NewLoader()
	if err := loader.Load("config/config.yaml"); err != nil {
		return nil, err
	}

	appCfg, err := loader.App()
	if err != nil {
		return nil, err
	}

	err = appCfg.Validate()
	if err != nil {
		return nil, err
	}

	databaseCfg, err := loader.Database()
	if err != nil {
		return nil, err
	}
	return &Config{
		AppCfg:      appCfg,
		DatabaseCfg: databaseCfg,
		Loader:      loader,
	}, nil
}
