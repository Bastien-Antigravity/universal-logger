package config

import (
	distributed_config "github.com/Bastien-Antigravity/distributed-config"
)

// DistConfig wraps the distributed configuration library.
type DistConfig struct {
	*distributed_config.Config
}

// -------------------------------------------------------------------------

// NewDistributedConfig initializes a new configuration service.
func NewDistributedConfig(profile string) *DistConfig {
	cfg := distributed_config.New(profile)
	return &DistConfig{
		Config: cfg,
	}
}

// -------------------------------------------------------------------------

// SetConfig updates a configuration value for a given section and key.
// Note: This specifically updates the in-memory configuration (LiveConfig).
// Subsystems monitoring updates via OnConfigUpdate will be notified.
func (s *DistConfig) SetConfig(section, key, value string) error {
	return s.Config.SetSingle(section, key, value)
}

// -------------------------------------------------------------------------

// Get returns a configuration value for a given section and key.
func (s *DistConfig) GetConfig(section, key string) string {
	return s.Config.Get(section, key)
}

// -------------------------------------------------------------------------

// OnConfigUpdate registers a callback for configuration updates.
func (s *DistConfig) OnConfigUpdate(fn func(map[string]map[string]string)) {
	s.Config.OnLiveConfUpdate(fn)
}
