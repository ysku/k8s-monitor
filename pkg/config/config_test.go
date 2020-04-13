package config

import "testing"

func TestValidate(t *testing.T) {
	config := NewConfig()
	config.EnableDeployment = false
	config.EnableJob = false
	config.EnablePersistentVolume = false
	config.EnablePod = false
	config.EnableService = false

	if err := config.Validate(); err == nil {
		t.Error("expected to get error, but not")
	}
}
