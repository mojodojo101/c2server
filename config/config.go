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

type ResourceConfig struct {
	TargetsPath string
	BeaconsPath string
}

func GetResourceConfig() *ResourceConfig {
	return &ResourceConfig{
		TargetsPath: "/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/",
		BeaconsPath: "/root/go/src/github.com/mojodojo101/c2server/internal_resources/beacons/",
	}
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
