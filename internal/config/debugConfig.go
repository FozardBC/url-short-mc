package config

import "time"

func DebugConfig() *Config {
	return &Config{
		Log: "debug",
		Storage: storage{
			Path:     "localhost:5432",
			Name:     "url-short-db",
			Username: "postgres",
			Password: "qwerty",
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
