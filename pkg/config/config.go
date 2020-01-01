package config

type Config struct {
	EnablePod     bool
	EnableService bool
}

func NewConfig() *Config {
	return &Config{}
}
