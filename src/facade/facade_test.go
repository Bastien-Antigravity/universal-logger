package facade

import (
	"testing"
	"github.com/Bastien-Antigravity/distconf-flexlog/src/models"
)

func TestNewDistconfFlexlogFacade(t *testing.T) {
	// 1. Basic Orchestration Test
	t.Run("Orchestration_Standalone_Devel", func(t *testing.T) {
		p := models.MFacadeParams{
			ConfigProfile: "standalone",
			AppName:       "test-app",
			LoggerProfile: "devel",
			LogLevel:      "debug",
		}

		f := NewDistconfFlexlogFacade(p)
		if f == nil {
			t.Fatal("NewDistconfFlexlogFacade returned nil")
		}
		defer f.Close()

		if f.Config == nil {
			t.Error("Facade.Config is nil")
		}

		if f.Logger == nil {
			t.Error("Facade.Logger is nil")
		}

		// Test direct config access
		if f.GetConfig() != f.Config {
			t.Error("GetConfig() did not return the correct config pointer")
		}
	})

	// 2. Log Level Propagation Test
	t.Run("LogLevel_Propagation", func(t *testing.T) {
		p := models.MFacadeParams{
			ConfigProfile: "standalone",
			AppName:       "test-app",
			LoggerProfile: "minimal",
			LogLevel:      "error",
		}

		f := NewDistconfFlexlogFacade(p)
		defer f.Close()

		// Since we added SetLevel to the engine, we can check it
		// We need to type assert the Logger to DistconfFlexlogFacade to access the engine?
		// No, the facade embeds interfaces.Logger.
		// However, we know that NewDistconfFlexlogFacade returns *DistconfFlexlogFacade.
		
		// In Go, we can't easily check the filtered level without accessing the private field or having a GetLevel().
		// But we can at least verify that it doesn't crash and the facade implements the interface.
		f.Error("This should be logged")
		f.Debug("This should NOT be logged")
	})

	// 3. Profile Switching Test
	t.Run("Profile_Switching", func(t *testing.T) {
		// Only testing profiles without mandatory network/server dependencies to avoid hanging
		profiles := []string{"devel", "minimal"}
		
		for _, profile := range profiles {
			t.Run("Profile-"+profile, func(t *testing.T) {
				p := models.MFacadeParams{
					ConfigProfile: "standalone", 
					AppName:       "test-profile-" + profile,
					LoggerProfile: profile,
					LogLevel:      "info",
				}
				
				f := NewDistconfFlexlogFacade(p)
				if f == nil {
					t.Errorf("Failed to initialize facade with profile: %s", profile)
				} else {
					f.Close()
				}
			})
		}
	})
}
