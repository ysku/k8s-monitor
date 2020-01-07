package config

import "golang.org/x/xerrors"

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

func (c *Config) Validate() error {
	if !(c.EnableDeployment || c.EnableJob || c.EnablePersistentVolume || c.EnablePod || c.EnableService) {
		return xerrors.New("at least one should be enabled!!")
	}
	return nil
}
