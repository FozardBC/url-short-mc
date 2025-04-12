package config

import "time"

func DebugConfig() *Config {
	return &Config{
		Log: "debug",
		Storage: storage{
			Path:     "localhost:8080",
			Name:     "url-s",
			Username: "postgres",
			Password: "Azer6789",
		},
		Server: http_server{
			Host:       "localhost",
			Port:       "8080",
			Timeout:    time.Duration(time.Second * 5),
			IdleTimout: time.Duration(time.Second * 60),
		},
		Rabbit: rabbit{
			Addr:     "",
			Login:    "",
			Password: "",
		},
	}
}
