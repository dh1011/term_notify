package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lee/term_notify/internal/notifier"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [command] [args...]",
	Short: "Run a command and notify when it finishes",
	Long: `Executes the given command, waits for it to complete, then sends
a push notification with the result, exit code, and duration.

Examples:
  tn run npm run build
  tn run ping -n 5 127.0.0.1
  tn -t my-builds run make -j8`,
	DisableFlagParsing: true,
	RunE:               runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified — usage: tn run <command> [args...]")
	}

	// Build the command
	name := args[0]
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	// On Windows, wrap in cmd /C for shell built-ins; on Unix, use sh -c
	var proc *exec.Cmd
	if runtime.GOOS == "windows" {
		allArgs := append([]string{"/C", name}, cmdArgs...)
		proc = exec.Command("cmd", allArgs...)
	} else {
		shellCmd := strings.Join(args, " ")
		proc = exec.Command("sh", "-c", shellCmd)
	}

	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	displayCmd := strings.Join(args, " ")
	fmt.Fprintf(os.Stderr, "tn: running %q\n", displayCmd)

	start := time.Now()
	err := proc.Run()
	elapsed := time.Since(start)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return fmt.Errorf("failed to run command: %w", err)
		}
	}

	// Build notification
	duration := formatDuration(elapsed)
	var title, body, tags string

	if exitCode == 0 {
		title = "✅ Command Succeeded"
		body = fmt.Sprintf("%s\nCompleted in %s", displayCmd, duration)
		tags = "white_check_mark"
	} else {
		title = "❌ Command Failed"
		body = fmt.Sprintf("%s\nFailed in %s (exit code %d)", displayCmd, duration, exitCode)
		tags = "x"
	}

	// Merge user tags
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
	} else {
		fmt.Fprintf(os.Stderr, "tn: notification sent → %s/%s\n", cfg.Server, cfg.Topic)
	}

	// Exit with the same code as the child process
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	if d < time.Hour {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	h := int(d.Hours())
	m = m % 60
	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}
