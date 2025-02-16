package app

type config struct {
	server  serverConfig
	storage storageConfig
	logger  loggerConfig
	crypto  cryptoConfig
}

type serverConfig struct {
	Addr string
}

type storageConfig struct {
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
		storageConfig{
			DatabaseURI: "",
		},
		loggerConfig{},
		cryptoConfig{
			HashKey:        "",
			PrivateKeyPath: "",
		},
	}, nil
}
