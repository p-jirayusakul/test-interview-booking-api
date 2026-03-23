package config

type Config struct {
	AppCfg      AppConfig
	DatabaseCfg DatabaseConfig
	Loader      *Loader
}

func InitConfig() (Config, error) {
	loader := NewLoader()
	if err := loader.Load("config/config.yaml"); err != nil {
		return Config{}, err
	}

	appCfg, err := loader.App()
	if err != nil {
		return Config{}, err
	}

	err = appCfg.Validate()
	if err != nil {
		return Config{}, err
	}

	databaseCfg, err := loader.Database()
	if err != nil {
		return Config{}, err
	}
	return Config{
		AppCfg:      appCfg,
		DatabaseCfg: databaseCfg,
		Loader:      loader,
	}, nil
}
