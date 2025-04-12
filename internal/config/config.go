package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Log     string      `mapstructure:"log"`
	Storage storage     `mapstructure:"storage"`
	Server  http_server `mapstructure:"http_server"`
	Rabbit  rabbit      `mapstructure:"rabbit"`
}

type rabbit struct {
	Addr     string `mapstructure:"addr"`
	Login    string `mapstructure:"login"`
	Password string `mapstructure:"password"`
}
type storage struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type http_server struct {
	Host       string        `mapstructure:"host"`
	Port       string        `mapstructure:"port"`
	Timeout    time.Duration `mapstructure:"timeout"`
	IdleTimout time.Duration `mapstructure:"idle_timeout"`
}

func MustReadConfig() *Config {
	//base settings
	viper := viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	// set defaults

	viper.SetDefault("http_server.host", "localhost")
	viper.SetDefault("http_server.port", 8080)

	// read cfg
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("can't read config: %s", err.Error())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicf("can't serialize config: %s", err.Error())
	}

	log.Print("Config is loaded sucesfully")

	return &cfg
}
