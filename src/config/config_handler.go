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

// SetConfig a configuration value for a given section and key.
// Note: This specifically updates the in-memory configuration (MemConfig).
// Subsystems monitoring updates via OnMemConfUpdate will be notified.
func (s *DistConfig) SetConfig(section, key, value string) {
	s.Config.Set(section, key, value)
}

// -------------------------------------------------------------------------

// Get returns a configuration value for a given section and key.
func (s *DistConfig) GetConfig(section, key string) string {
	return s.Config.Get(section, key)
}

// -------------------------------------------------------------------------

// OnConfigUpdate registers a callback for configuration updates.
func (s *DistConfig) OnConfigUpdate(fn func(map[string]map[string]string)) {
	s.Config.OnMemConfUpdate(fn)
}
