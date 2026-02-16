package cmd

import (
	"fmt"
	"os"

	"github.com/lee/term_notify/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg    *config.Config
	cfgErr error

	// Flag overrides
	flagServer   string
	flagTopic    string
	flagPriority string
	flagTags     string
)

var rootCmd = &cobra.Command{
	Use:   "tn",
	Short: "term_notify — get notified when terminal commands finish",
	Long: `term_notify (tn) sends push notifications via ntfy when your
terminal commands complete. Wrap a command, watch a PID, or
send a quick notification.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&flagServer, "server", "s", "", "ntfy server (default: ntfy.sh)")
	rootCmd.PersistentFlags().StringVarP(&flagTopic, "topic", "t", "", "ntfy topic name")
	rootCmd.PersistentFlags().StringVarP(&flagPriority, "priority", "p", "", "notification priority (min, low, default, high, max)")
	rootCmd.PersistentFlags().StringVar(&flagTags, "tags", "", "comma-separated tags/emojis")
}

func initConfig() {
	cfg, cfgErr = config.Load()
	if cfgErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load config: %v\n", cfgErr)
		cfg = config.DefaultConfig()
	}

	// CLI flags override config values
	if flagServer != "" {
		cfg.Server = flagServer
	}
	if flagTopic != "" {
		cfg.Topic = flagTopic
	}
	if flagPriority != "" {
		cfg.Priority = flagPriority
	}
}

// getEffectiveTags returns the tags to use — flag takes precedence.
func getEffectiveTags() string {
	return flagTags
}
