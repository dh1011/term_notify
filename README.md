# term_notify (tn)

A lightweight CLI tool that sends push notifications via [ntfy](https://ntfy.sh) when your terminal commands finish. Works on **Windows**, **macOS**, and **Linux** with both **PowerShell** and **Bash/Zsh**.

## Quick Start

```bash
# 1. Build
go build -o tn.exe .   # Windows
go build -o tn .        # Linux/macOS

# 2. Set your ntfy topic (one-time setup)
tn config --topic my-term-alerts

# 3. Subscribe to your topic
#    - Phone: Install the ntfy app, subscribe to "my-term-alerts"
#    - Browser: Visit https://ntfy.sh/my-term-alerts

# 4. Use it!
tn run npm run build
```

## Commands

### `tn run <command> [args...]`

Wraps a command — runs it, then sends a notification with the result.

```bash
tn run make -j8
tn run python train.py --epochs 100
tn run docker build -t myapp .
```

The notification includes:
- ✅/❌ Success or failure
- Command name
- Duration
- Exit code (on failure)

### `tn pid <process-id>`

Watches an already-running process and notifies when it exits.

```bash
# In another terminal, find the PID of a running process
tn pid 12345
```

### `tn notify <message>`

Sends a one-shot notification. Great for chaining.

```powershell
# PowerShell
npm run build; tn notify "Build finished"

# Bash
npm run build && tn notify "Build succeeded" || tn notify "Build failed"
```

### `tn config`

View or update your configuration.

```bash
tn config --topic my-alerts      # Set topic (required)
tn config --server ntfy.sh       # Set server (default: ntfy.sh)
tn config --priority high        # Set priority (min/low/default/high/max)
tn config --token tk_xxx         # Set auth token (for private servers)
tn config                        # Show current config
```

## Configuration

Settings are stored in a YAML config file:

| OS      | Path                                      |
|---------|-------------------------------------------|
| Windows | `%APPDATA%\term_notify\config.yaml`       |
| Linux   | `~/.config/term_notify/config.yaml`       |
| macOS   | `~/.config/term_notify/config.yaml`       |

```yaml
server: ntfy.sh
topic: my-term-alerts
priority: default
token: ""
```

### Environment Variables

Environment variables override config file values:

| Variable      | Description       |
|---------------|-------------------|
| `TN_SERVER`   | ntfy server URL   |
| `TN_TOPIC`    | ntfy topic name   |
| `TN_PRIORITY` | Default priority  |
| `TN_TOKEN`    | Auth token        |

### CLI Flags

Flags override both config file and environment variables:

```bash
tn run --topic urgent-builds --priority high make build
```

**Precedence:** CLI flags > env vars > config file > defaults

## Shell Integration

### PowerShell

Add to your `$PROFILE` (`~\Documents\PowerShell\Microsoft.PowerShell_profile.ps1`):

```powershell
. "E:\Development\term_notify\shell\tn.ps1"
```

This gives you:
- `tnr` — alias for `tn run` (e.g., `tnr npm run build`)
- `tnd` — notify about the last command's result (e.g., `npm run build; tnd`)

### Bash / Zsh

Add to your `.bashrc` or `.zshrc`:

```bash
source /path/to/term_notify/shell/tn.bash
```

Same aliases: `tnr` and `tnd`.

The bash script also includes a commented-out **auto-notify** hook that automatically notifies for any command running longer than 30 seconds.

## Self-Hosted ntfy

To use your own ntfy server:

```bash
tn config --server ntfy.example.com
tn config --token tk_your_token_here
```

## Building from Source

```bash
git clone https://github.com/lee/term_notify.git
cd term_notify
go build -o tn.exe .    # Windows
go build -o tn .         # Linux/macOS
```

### Cross-compile

```bash
GOOS=linux GOARCH=amd64 go build -o tn .
GOOS=darwin GOARCH=arm64 go build -o tn .
GOOS=windows GOARCH=amd64 go build -o tn.exe .
```

## License

MIT
