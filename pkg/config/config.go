package config

type Config struct {
	EnableDeployment       bool
	EnableJob              bool
	EnablePersistentVolume bool
	EnablePod              bool
	EnableService          bool
}

func NewConfig() *Config {
	return &Config{}
}
