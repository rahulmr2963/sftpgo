package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sftpgo/sdk/plugin/notifier"
)

const (
	pluginName = "sftpgo-ssh-key-logger"
	version    = "1.0.0"
)

type sshKeyLoggerPlugin struct {
	logger   hclog.Logger
	logFile  string
}

// NewSSHKeyLoggerPlugin creates a new SSH key logger plugin instance
func NewSSHKeyLoggerPlugin() *sshKeyLoggerPlugin {
	return &sshKeyLoggerPlugin{
		logger:  hclog.New(&hclog.LoggerOptions{Name: pluginName}),
		logFile: getLogFilePath(),
	}
}

func getLogFilePath() string {
	// Get log file path from environment or use default
	logPath := os.Getenv("SSH_KEY_LOG_PATH")
	if logPath == "" {
		logPath = "/tmp/sftpgo-ssh-keys.log"
	}
	return logPath
}

// NotifyFsEvent implements the notifier interface for filesystem events
func (p *sshKeyLoggerPlugin) NotifyFsEvent(event *notifier.FsEvent) error {
	// Log FS events for debugging
	p.logger.Debug("Received FS event", 
		"action", event.Action, 
		"username", event.Username,
		"path", event.Path)
	
	// We're not specifically interested in FS events for SSH key logging
	return nil
}

// NotifyProviderEvent implements the notifier interface for provider events
func (p *sshKeyLoggerPlugin) NotifyProviderEvent(event *notifier.ProviderEvent) error {
	// Log provider events for debugging
	p.logger.Debug("Received provider event", 
		"action", event.Action, 
		"username", event.Username,
		"object_type", event.ObjectType,
		"object_name", event.ObjectName)

	// Check if this is a user update that might include SSH key changes
	if event.Action == "update" && event.ObjectType == "user" {
		return p.logSSHKeyEvent(event)
	}
	
	return nil
}

// NotifyLogEvent implements the notifier interface for log events
func (p *sshKeyLoggerPlugin) NotifyLogEvent(event *notifier.LogEvent) error {
	// Log events for debugging
	p.logger.Debug("Received log event", 
		"event", event.Event, 
		"username", event.Username,
		"ip", event.IP)

	// Check if this might be related to SSH key operations
	// For example, if we added specific logging for SSH key generation
	if p.isSSHKeyLogEvent(event) {
		return p.logSSHKeyFromLogEvent(event)
	}
	
	return nil
}

func (p *sshKeyLoggerPlugin) isSSHKeyLogEvent(event *notifier.LogEvent) bool {
	// This would check if the log event is related to SSH key generation
	// We could look for specific log messages that contain "SSH key"
	return false  // Placeholder - in practice you'd implement real detection
}

func (p *sshKeyLoggerPlugin) logSSHKeyEvent(event *notifier.ProviderEvent) error {
	// Open log file for appending
	file, err := os.OpenFile(p.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		p.logger.Error("Failed to open log file", "error", err, "file", p.logFile)
		return err
	}
	defer file.Close()

	// Create log entry
	timestamp := time.Unix(event.Timestamp, 0).Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] User Update Event - User: %s, Action: %s, Object: %s, IP: %s\n",
		timestamp, event.Username, event.Action, event.ObjectName, event.IP)

	// Write to log file
	if _, err := file.WriteString(logEntry); err != nil {
		p.logger.Error("Failed to write to log file", "error", err)
		return err
	}

	p.logger.Info("User update event logged (potential SSH key change)", "user", event.Username, "file", p.logFile)
	return nil
}

func (p *sshKeyLoggerPlugin) logSSHKeyFromLogEvent(event *notifier.LogEvent) error {
	// Similar logging for log events
	file, err := os.OpenFile(p.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		p.logger.Error("Failed to open log file", "error", err, "file", p.logFile)
		return err
	}
	defer file.Close()

	timestamp := time.Unix(event.Timestamp, 0).Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] SSH Key Log Event - User: %s, Event: %d, IP: %s\n",
		timestamp, event.Username, event.Event, event.IP)

	if _, err := file.WriteString(logEntry); err != nil {
		p.logger.Error("Failed to write to log file", "error", err)
		return err
	}

	p.logger.Info("SSH key log event recorded", "user", event.Username, "file", p.logFile)
	return nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   pluginName,
		Output: os.Stderr,
		Level:  hclog.Info,
	})

	logger.Info("Starting SSH Key Logger Plugin", "version", version)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: notifier.Handshake,
		Plugins: map[string]plugin.Plugin{
			notifier.PluginName: &notifier.Plugin{Impl: NewSSHKeyLoggerPlugin()},
		},
		Logger: logger,
	})
}