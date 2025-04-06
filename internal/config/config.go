package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
	Metrics []MetricConfig `yaml:"metrics" json:"metrics"`
}

type ServerConfig struct {
	Addres string `yaml:"host" json:"host" env-default:"127.0.0.1"`
	Port string `yaml:"port" json:"port" env-default:"8152"`
}

type MetricConfig struct {
	Name string `yaml:"name" json:"name"`
	Options MetricOptions `yaml:"options" json:"options"`
}

type MetricOptions struct {
	Name string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

func MustLoad() *Config {
	path := loadPath()
	if path == "" {
		panic("Can`t read config file")
	}

	return loadConfig(path)
}

func loadPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {	
		path = ".././config/config.yml"
	}

	return path
}

func loadConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	
	return &cfg
}