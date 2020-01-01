package config

type Config struct {
	EnableDeployment bool
	EnablePod        bool
	EnableService    bool
}

func NewConfig() *Config {
	return &Config{}
}
