package config

import(

)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Addres string
	Port string
	PathToMetrics string
}

func New() *Config {
	return &Config{}
}