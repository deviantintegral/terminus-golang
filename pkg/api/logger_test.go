package api

import (
	"bytes"
	"strings"
	"testing"
)

func TestRedactSensitiveData(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "redact machine_token in request body",
			input:    `{"machine_token": "abcdefghijklmnopqrstuvwxyz123456", "client": "terminus-golang"}`,
			expected: `{"machine_token": "[REDACTED]", "client": "terminus-golang"}`,
		},
		{
			name:     "redact session in response body",
			input:    `{"session": "abcdefghijklmnopqrstuvwxyz123456", "user_id": "12345", "expires_at": 1234567890}`,
			expected: `{"session": "[REDACTED]", "user_id": "12345", "expires_at": 1234567890}`,
		},
		{
			name:     "redact session_token in body",
			input:    `{"session_token": "abcdefghijklmnopqrstuvwxyz123456", "user_id": "12345"}`,
			expected: `{"session_token": "[REDACTED]", "user_id": "12345"}`,
		},
		{
			name:     "do not redact short tokens",
			input:    `{"session": "short", "machine_token": "tiny"}`,
			expected: `{"session": "short", "machine_token": "tiny"}`,
		},
		{
			name:     "handle whitespace variations",
			input:    `{"machine_token":"abcdefghijklmnopqrstuvwxyz123456"}`,
			expected: `{"machine_token": "[REDACTED]"}`,
		},
		{
			name:     "handle multiple tokens in same body",
			input:    `{"machine_token": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "session": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}`,
			expected: `{"machine_token": "[REDACTED]", "session": "[REDACTED]"}`,
		},
		{
			name:     "leave non-sensitive fields unchanged",
			input:    `{"client": "terminus-golang", "user_id": "12345"}`,
			expected: `{"client": "terminus-golang", "user_id": "12345"}`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "redact token with escaped quote",
			input:    `{"machine_token": "abcdefghij\"klmnopqrstuvwxyz"}`,
			expected: `{"machine_token": "[REDACTED]"}`,
		},
		{
			name:     "redact token with escaped backslash",
			input:    `{"session": "abcdefghij\\klmnopqrstuvwxyz"}`,
			expected: `{"session": "[REDACTED]"}`,
		},
		{
			name:     "redact token with multiple escapes",
			input:    `{"machine_token": "abc\"def\\ghi\"jkl\\mnopqrst"}`,
			expected: `{"machine_token": "[REDACTED]"}`,
		},
		{
			name:     "redact token with unicode escape",
			input:    `{"session": "abcdefghij\u0041klmnopqrstuvwxyz"}`,
			expected: `{"session": "[REDACTED]"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactSensitiveData(tt.input)
			if result != tt.expected {
				t.Errorf("redactSensitiveData() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestLogHTTPRequestRedactsBody(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter(VerbosityTrace, &buf)

	body := `{"machine_token": "abcdefghijklmnopqrstuvwxyz123456", "client": "terminus-golang"}`
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	logger.LogHTTPRequest("POST", "https://example.com/api/test", headers, body)

	output := buf.String()

	// Should contain redacted token
	if !strings.Contains(output, `"machine_token": "[REDACTED]"`) {
		t.Errorf("expected redacted machine_token in output, got: %s", output)
	}

	// Should NOT contain original token
	if strings.Contains(output, "abcdefghijklmnopqrstuvwxyz123456") {
		t.Errorf("output should not contain original token, got: %s", output)
	}
}

func TestLogHTTPResponseRedactsBody(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter(VerbosityTrace, &buf)

	body := `{"session": "abcdefghijklmnopqrstuvwxyz123456", "user_id": "12345", "expires_at": 1234567890}`
	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}

	logger.LogHTTPResponse(200, "OK", headers, body)

	output := buf.String()

	// Should contain redacted token
	if !strings.Contains(output, `"session": "[REDACTED]"`) {
		t.Errorf("expected redacted session in output, got: %s", output)
	}

	// Should NOT contain original token
	if strings.Contains(output, "abcdefghijklmnopqrstuvwxyz123456") {
		t.Errorf("output should not contain original token, got: %s", output)
	}
}

func TestLogHTTPRequestRedactsAuthorizationHeader(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLoggerWithWriter(VerbosityTrace, &buf)

	headers := map[string][]string{
		"Authorization": {"Bearer my-secret-token"},
		"Content-Type":  {"application/json"},
	}

	logger.LogHTTPRequest("GET", "https://example.com/api/test", headers, "")

	output := buf.String()

	// Should contain redacted Authorization header
	if !strings.Contains(output, "Authorization: [REDACTED]") {
		t.Errorf("expected redacted Authorization header in output, got: %s", output)
	}

	// Should NOT contain original token
	if strings.Contains(output, "my-secret-token") {
		t.Errorf("output should not contain original token, got: %s", output)
	}
}

func TestLoggerVerbosityLevels(t *testing.T) {
	tests := []struct {
		name          string
		verbosity     VerbosityLevel
		method        string
		shouldContain []string
		shouldNotLog  bool
	}{
		{
			name:         "VerbosityNone logs nothing",
			verbosity:    VerbosityNone,
			method:       "debug",
			shouldNotLog: true,
		},
		{
			name:          "VerbosityInfo logs info",
			verbosity:     VerbosityInfo,
			method:        "info",
			shouldContain: []string{"[INFO]"},
		},
		{
			name:         "VerbosityInfo does not log debug",
			verbosity:    VerbosityInfo,
			method:       "debug",
			shouldNotLog: true,
		},
		{
			name:          "VerbosityDebug logs debug",
			verbosity:     VerbosityDebug,
			method:        "debug",
			shouldContain: []string{"[DEBUG]"},
		},
		{
			name:          "VerbosityTrace logs trace",
			verbosity:     VerbosityTrace,
			method:        "trace",
			shouldContain: []string{"[TRACE]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLoggerWithWriter(tt.verbosity, &buf)

			switch tt.method {
			case "debug":
				logger.Debug("test message")
			case "info":
				logger.Info("test message")
			case "trace":
				logger.Trace("test message")
			}

			output := buf.String()

			if tt.shouldNotLog {
				if output != "" {
					t.Errorf("expected no output, got: %s", output)
				}
				return
			}

			for _, expected := range tt.shouldContain {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestIsTraceEnabled(t *testing.T) {
	tests := []struct {
		verbosity VerbosityLevel
		expected  bool
	}{
		{VerbosityNone, false},
		{VerbosityInfo, false},
		{VerbosityDebug, false},
		{VerbosityTrace, true},
	}

	for _, tt := range tests {
		logger := NewLogger(tt.verbosity)
		if logger.IsTraceEnabled() != tt.expected {
			t.Errorf("IsTraceEnabled() for verbosity %d = %v, expected %v",
				tt.verbosity, logger.IsTraceEnabled(), tt.expected)
		}
	}
}

func TestAsHTTPLogger(t *testing.T) {
	// Test with nil
	httpLogger, ok := AsHTTPLogger(nil)
	if ok || httpLogger != nil {
		t.Error("expected AsHTTPLogger(nil) to return nil, false")
	}

	// Test with DefaultLogger (implements HTTPLogger)
	defaultLogger := NewLogger(VerbosityTrace)
	httpLogger, ok = AsHTTPLogger(defaultLogger)
	if !ok || httpLogger == nil {
		t.Error("expected AsHTTPLogger(DefaultLogger) to return logger, true")
	}
}
