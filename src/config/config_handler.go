package config

import (
	distributed_config "github.com/Bastien-Antigravity/distributed-config/src/facade"
)

// DistConfig wraps the distributed configuration library.
type DistConfig struct {
	*distributed_config.Config
}

// -------------------------------------------------------------------------
// NewDistributedConfig initializes a new configuration service.
func NewDistributedConfig(profile string) *DistConfig {
	cfg := distributed_config.NewConfig(profile)
	return &DistConfig{
		Config: cfg,
	}
}

// -------------------------------------------------------------------------
// Get returns a configuration value for a given section and key.
func (s *DistConfig) Get(section, key string) string {
	return s.Config.Get(section, key)
}

// -------------------------------------------------------------------------
// OnMemConfUpdate registers a callback for configuration updates.
func (s *DistConfig) OnMemConfUpdate(fn func(map[string]map[string]string)) {
	s.Config.OnMemConfUpdate(fn)
}
