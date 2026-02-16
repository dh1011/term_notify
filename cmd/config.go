package cmd

import (
	"fmt"

	"github.com/lee/term_notify/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update configuration",
	Long: `View or set term_notify configuration values.
Settings are stored in a YAML file in your config directory.

Examples:
  tn config --topic my-alerts        # Set the ntfy topic
  tn config --server ntfy.example.com  # Use a self-hosted server
  tn config                           # Show current config`,
	RunE: runConfig,
}

var configTopic string
var configServer string
var configPriority string
var configToken string

func init() {
	configCmd.Flags().StringVar(&configTopic, "topic", "", "set ntfy topic")
	configCmd.Flags().StringVar(&configServer, "server", "", "set ntfy server")
	configCmd.Flags().StringVar(&configPriority, "priority", "", "set default priority")
	configCmd.Flags().StringVar(&configToken, "token", "", "set auth token")
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	changed := false

	if configTopic != "" {
		cfg.Topic = configTopic
		changed = true
	}
	if configServer != "" {
		cfg.Server = configServer
		changed = true
	}
	if configPriority != "" {
		cfg.Priority = configPriority
		changed = true
	}
	if configToken != "" {
		cfg.Token = configToken
		changed = true
	}

	if changed {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		path, _ := config.ConfigPath()
		fmt.Printf("Config saved to %s\n", path)
	}

	// Always display current config
	fmt.Println()
	fmt.Printf("  server:   %s\n", cfg.Server)
	fmt.Printf("  topic:    %s\n", displayValue(cfg.Topic))
	fmt.Printf("  priority: %s\n", cfg.Priority)
	fmt.Printf("  token:    %s\n", maskToken(cfg.Token))
	fmt.Println()

	if cfg.Topic == "" {
		fmt.Println("⚠️  No topic set. Run: tn config --topic <your-topic>")
	}

	return nil
}

func displayValue(v string) string {
	if v == "" {
		return "(not set)"
	}
	return v
}

func maskToken(t string) string {
	if t == "" {
		return "(not set)"
	}
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "****" + t[len(t)-4:]
}
