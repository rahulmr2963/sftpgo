# SFTPGo SSH Key Logger Plugin

This is an example plugin for SFTPGo that demonstrates how to create a notifier plugin to log SSH key related events.

## Features

- Logs SSH key related events to a file
- Configurable log file path via environment variable
- Tracks user updates that may include SSH key changes
- Demonstrates SFTPGo plugin development patterns

## Building

1. Initialize the Go module and download dependencies:
```bash
cd /tmp/sftpgo-ssh-key-logger-plugin
go mod tidy
```

2. Build the plugin:
```bash
go build -o sftpgo-ssh-key-logger-plugin
```

## Configuration

Add this plugin to your SFTPGo configuration:

```json
{
  "plugins": [
    {
      "type": "notifier",
      "cmd": "/path/to/sftpgo-ssh-key-logger-plugin",
      "notifier_options": {
        "provider_events": ["update"],
        "provider_objects": ["user"]
      }
    }
  ]
}
```

## Environment Variables

- `SSH_KEY_LOG_PATH`: Path to the log file (default: `/tmp/sftpgo-ssh-keys.log`)

## Usage

Once configured and SFTPGo is restarted, the plugin will:

1. Monitor user update events
2. Log SSH key related activities to the specified log file
3. Provide audit trail for SSH key management

## Log Format

```
[2024-09-06T10:30:45Z] SSH Key Event - User: john_doe, Event: update, Status: 1, IP: 192.168.1.100
```

## Extending the Plugin

This plugin can be extended to:

- Parse user data to detect actual SSH key changes
- Send notifications to external systems (email, Slack, etc.)
- Implement key rotation policies
- Integrate with key management systems
- Validate SSH key formats and security policies

## Notes

- This is a demonstration plugin showing the concepts
- In production, you may want to add more sophisticated event detection
- Consider implementing proper error handling and retry logic
- The plugin runs as a separate process and communicates with SFTPGo via RPC