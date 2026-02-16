package notifier

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message represents a notification to be sent.
type Message struct {
	Server   string
	Topic    string
	Title    string
	Body     string
	Priority string
	Tags     string
	Token    string
}

// Send publishes a notification message to the ntfy server.
func Send(msg *Message) error {
	if msg.Topic == "" {
		return fmt.Errorf("topic is required — run 'tn config --topic <name>' or set TN_TOPIC")
	}

	server := msg.Server
	if server == "" {
		server = "ntfy.sh"
	}

	// Ensure server has a scheme
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "https://" + server
	}

	url := fmt.Sprintf("%s/%s", strings.TrimRight(server, "/"), msg.Topic)

	req, err := http.NewRequest("POST", url, strings.NewReader(msg.Body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if msg.Title != "" {
		req.Header.Set("Title", msg.Title)
	}
	if msg.Priority != "" && msg.Priority != "default" {
		req.Header.Set("Priority", msg.Priority)
	}
	if msg.Tags != "" {
		req.Header.Set("Tags", msg.Tags)
	}
	if msg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+msg.Token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ntfy server returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
