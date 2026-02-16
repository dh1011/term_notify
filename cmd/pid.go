package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lee/term_notify/internal/notifier"
	"github.com/lee/term_notify/internal/process"
	"github.com/spf13/cobra"
)

var pidCmd = &cobra.Command{
	Use:   "pid <process-id>",
	Short: "Watch a running process and notify when it exits",
	Long: `Monitors an already-running process by its PID. When the process
exits, a push notification is sent.

Examples:
  tn pid 12345
  tn pid --topic builds 12345`,
	Args: cobra.ExactArgs(1),
	RunE: runPid,
}

func init() {
	rootCmd.AddCommand(pidCmd)
}

func runPid(cmd *cobra.Command, args []string) error {
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PID %q: %w", args[0], err)
	}

	fmt.Fprintf(os.Stderr, "tn: watching PID %d…\n", pid)

	elapsed, err := process.WaitForPID(pid)
	if err != nil {
		return fmt.Errorf("watching process: %w", err)
	}

	duration := formatDuration(elapsed)
	title := "🏁 Process Exited"
	body := fmt.Sprintf("PID %d exited after %s", pid, duration)
	tags := "checkered_flag"

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

	if notifyErr := notifier.Send(msg); notifyErr != nil {
		fmt.Fprintf(os.Stderr, "tn: notification failed: %v\n", notifyErr)
		return notifyErr
	}

	fmt.Fprintf(os.Stderr, "tn: PID %d finished — notification sent → %s/%s\n", pid, cfg.Server, cfg.Topic)
	return nil
}
