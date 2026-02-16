package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server != "ntfy.sh" {
		t.Errorf("DefaultConfig().Server = %q, want %q", cfg.Server, "ntfy.sh")
	}
	if cfg.Priority != "default" {
		t.Errorf("DefaultConfig().Priority = %q, want %q", cfg.Priority, "default")
	}
	if cfg.Topic != "" {
		t.Errorf("DefaultConfig().Topic = %q, want empty", cfg.Topic)
	}
	if cfg.Token != "" {
		t.Errorf("DefaultConfig().Token = %q, want empty", cfg.Token)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	t.Run("all env vars set", func(t *testing.T) {
		t.Setenv("TN_SERVER", "custom.ntfy.example.com")
		t.Setenv("TN_TOPIC", "test-topic")
		t.Setenv("TN_TOKEN", "secret-token-123")
		t.Setenv("TN_PRIORITY", "high")

		cfg := DefaultConfig()
		applyEnvOverrides(cfg)

		if cfg.Server != "custom.ntfy.example.com" {
			t.Errorf("Server = %q, want %q", cfg.Server, "custom.ntfy.example.com")
		}
		if cfg.Topic != "test-topic" {
			t.Errorf("Topic = %q, want %q", cfg.Topic, "test-topic")
		}
		if cfg.Token != "secret-token-123" {
			t.Errorf("Token = %q, want %q", cfg.Token, "secret-token-123")
		}
		if cfg.Priority != "high" {
			t.Errorf("Priority = %q, want %q", cfg.Priority, "high")
		}
	})

	t.Run("empty env vars keep defaults", func(t *testing.T) {
		// Ensure these are unset for this test
		t.Setenv("TN_SERVER", "")
		t.Setenv("TN_TOPIC", "")
		t.Setenv("TN_TOKEN", "")
		t.Setenv("TN_PRIORITY", "")

		cfg := DefaultConfig()
		applyEnvOverrides(cfg)

		if cfg.Server != "ntfy.sh" {
			t.Errorf("Server = %q, want %q", cfg.Server, "ntfy.sh")
		}
		if cfg.Priority != "default" {
			t.Errorf("Priority = %q, want %q", cfg.Priority, "default")
		}
	})

	t.Run("partial env override", func(t *testing.T) {
		t.Setenv("TN_SERVER", "")
		t.Setenv("TN_TOPIC", "partial-topic")
		t.Setenv("TN_TOKEN", "")
		t.Setenv("TN_PRIORITY", "")

		cfg := DefaultConfig()
		applyEnvOverrides(cfg)

		if cfg.Server != "ntfy.sh" {
			t.Errorf("Server = %q, want default %q", cfg.Server, "ntfy.sh")
		}
		if cfg.Topic != "partial-topic" {
			t.Errorf("Topic = %q, want %q", cfg.Topic, "partial-topic")
		}
	})
}

func TestConfigYAMLRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	original := &Config{
		Server:   "my-server.example.com",
		Topic:    "my-topic",
		Priority: "high",
		Token:    "my-secret-token",
	}

	// Marshal to YAML and write
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	// Read back and unmarshal
	loaded := &Config{}
	fileData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read temp config: %v", err)
	}
	if err := yaml.Unmarshal(fileData, loaded); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if loaded.Server != original.Server {
		t.Errorf("Server = %q, want %q", loaded.Server, original.Server)
	}
	if loaded.Topic != original.Topic {
		t.Errorf("Topic = %q, want %q", loaded.Topic, original.Topic)
	}
	if loaded.Priority != original.Priority {
		t.Errorf("Priority = %q, want %q", loaded.Priority, original.Priority)
	}
	if loaded.Token != original.Token {
		t.Errorf("Token = %q, want %q", loaded.Token, original.Token)
	}
}
