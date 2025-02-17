package app

type config struct {
	server     serverConfig
	repository repositoryConfig
	logger     loggerConfig
	crypto     cryptoConfig
}

type serverConfig struct {
	Addr string
}

type repositoryConfig struct {
	DatabaseURI string
}

type loggerConfig struct {
	//LogLevel zap.AtomicLevel
}

type cryptoConfig struct {
	HashKey        string
	PrivateKeyPath string
}

func newConfig() (*config, error) {
	return &config{
		serverConfig{
			Addr: "",
		},
		repositoryConfig{
			DatabaseURI: "",
		},
		loggerConfig{},
		cryptoConfig{
			HashKey:        "",
			PrivateKeyPath: "",
		},
	}, nil
}
