package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lee/term_notify/internal/notifier"
	"github.com/spf13/cobra"
)

var notifyTitle string

var notifyCmd = &cobra.Command{
	Use:   "notify <message>",
	Short: "Send an ad-hoc notification",
	Long: `Sends a one-shot push notification with the given message.
Useful for chaining with other commands.

Examples:
  tn notify "Build complete!"
  make build; tn notify "Build finished"
  tn notify --title "Deploy" "Deployed to production"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runNotify,
}

func init() {
	notifyCmd.Flags().StringVar(&notifyTitle, "title", "", "notification title")
	rootCmd.AddCommand(notifyCmd)
}

func runNotify(cmd *cobra.Command, args []string) error {
	body := strings.Join(args, " ")
	title := notifyTitle
	if title == "" {
		title = "📢 term_notify"
	}

	tags := "loudspeaker"
	if userTags := getEffectiveTags(); userTags != "" {
		tags = tags + "," + userTags
	}

	msg := &notifier.Message{
		Server:   cfg.Server,
		Topic:    cfg.Topic,
		Title:    title,
		Body:     body,
		Priority: cfg.Priority,
		Tags:     tags,
		Token:    cfg.Token,
	}

	if err := notifier.Send(msg); err != nil {
		fmt.Fprintf(os.Stderr, "tn: notification failed: %v\n", err)
		return err
	}

	fmt.Fprintf(os.Stderr, "tn: notification sent → %s/%s\n", cfg.Server, cfg.Topic)
	return nil
}
