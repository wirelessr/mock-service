package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"mock-service/internal/models"
)

// LoggerImpl implements the Logger interface
// It provides structured logging functionality for the mock service
type LoggerImpl struct{}

// NewLogger creates a new instance of Logger
func NewLogger() *LoggerImpl {
	return &LoggerImpl{}
}

// LogRequest logs incoming HTTP request details in JSON format
func (l *LoggerImpl) LogRequest(method, path string, params map[string]string) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     "INFO",
		"type":      "request",
		"method":    method,
		"path":      path,
		"params":    params,
	}

	l.writeLog(logEntry)
}

// LogResponse logs outgoing HTTP response details in JSON format
func (l *LoggerImpl) LogResponse(statusCode int, body interface{}) {
	logEntry := map[string]interface{}{
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"level":       "INFO",
		"type":        "response",
		"status_code": statusCode,
		"body":        body,
	}

	l.writeLog(logEntry)
}

// LogMatch logs when a rule is matched in JSON format
func (l *LoggerImpl) LogMatch(rule *models.MockRule) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     "INFO",
		"type":      "match",
		"message":   "Rule matched",
		"rule": map[string]interface{}{
			"path": rule.Path,
			"code": rule.Code,
		},
	}

	l.writeLog(logEntry)
}

// LogDefault logs when default response is used in JSON format
func (l *LoggerImpl) LogDefault() {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"level":     "INFO",
		"type":      "default",
		"message":   "No matching rule found, using default response",
	}

	l.writeLog(logEntry)
}

// writeLog writes the log entry to stdout in JSON format
func (l *LoggerImpl) writeLog(logEntry map[string]interface{}) {
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		// Fallback to simple error message if JSON marshaling fails
		fmt.Fprintf(os.Stdout, `{"timestamp":"%s","level":"ERROR","type":"log_error","message":"Failed to marshal log entry"}`+"\n",
			time.Now().UTC().Format(time.RFC3339))
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonData))
}
