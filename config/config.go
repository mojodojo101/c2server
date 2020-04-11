package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	DBName   string
	Host     string
	Port     int
	Username string
	Password string
	SSLMode  string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			DBName:   "c2db",
			Host:     "localhost",
			Port:     5432,
			Username: "c2admin",
			Password: "mojodojo101+",
			SSLMode:  "require",
		},
	}
}
