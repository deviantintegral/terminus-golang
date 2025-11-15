package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultHost is the default Pantheon API host
	DefaultHost = "terminus.pantheon.io"
	// DefaultPort is the default API port
	DefaultPort = 443
	// DefaultProtocol is the default protocol
	DefaultProtocol = "https"
	// DefaultTimeout is the default timeout in seconds
	DefaultTimeout = 86400
	// DefaultDateFormat is the default date format
	DefaultDateFormat = "2006-01-02 15:04:05"
)

// Config represents the application configuration
type Config struct {
	// API settings
	Host     string
	Port     int
	Protocol string
	Timeout  int

	// Paths
	HomeDir      string
	CacheDir     string
	PluginsDir   string
	TokensDir    string
	SessionFile  string
	ConfigFile   string

	// Display settings
	DateFormat string

	// Other settings
	values map[string]interface{}
}

// New creates a new configuration with defaults
func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	terminusDir := filepath.Join(homeDir, ".terminus")
	cacheDir := filepath.Join(terminusDir, "cache")
	tokensDir := filepath.Join(cacheDir, "tokens")
	pluginsDir := filepath.Join(terminusDir, "plugins-3.x")

	cfg := &Config{
		Host:        DefaultHost,
		Port:        DefaultPort,
		Protocol:    DefaultProtocol,
		Timeout:     DefaultTimeout,
		DateFormat:  DefaultDateFormat,
		HomeDir:     homeDir,
		CacheDir:    cacheDir,
		PluginsDir:  pluginsDir,
		TokensDir:   tokensDir,
		SessionFile: filepath.Join(cacheDir, "session"),
		ConfigFile:  filepath.Join(terminusDir, "config.yml"),
		values:      make(map[string]interface{}),
	}

	// Load configuration layers
	if err := cfg.load(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// load loads configuration from all sources in priority order
func (c *Config) load() error {
	// 1. Load from user config file if it exists
	if err := c.loadFromFile(c.ConfigFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// 2. Load from .env file in current directory if it exists
	if err := c.loadFromEnvFile(".env"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	// 3. Load from environment variables (highest priority)
	c.loadFromEnvironment()

	// Ensure directories exist
	if err := c.ensureDirectories(); err != nil {
		return err
	}

	return nil
}

// loadFromFile loads configuration from a YAML file
func (c *Config) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var values map[string]interface{}
	if err := yaml.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	for key, value := range values {
		c.setValue(key, value)
	}

	return nil
}

// loadFromEnvFile loads environment variables from a .env file
func (c *Config) loadFromEnvFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		if strings.HasPrefix(key, "TERMINUS_") {
			c.setValue(key, value)
		}
	}

	return nil
}

// loadFromEnvironment loads configuration from environment variables
func (c *Config) loadFromEnvironment() {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		if strings.HasPrefix(key, "TERMINUS_") {
			c.setValue(key, value)
		}
	}
}

// setValue sets a configuration value, handling special keys
func (c *Config) setValue(key string, value interface{}) {
	// Normalize key to uppercase with TERMINUS_ prefix
	normalizedKey := strings.ToUpper(key)
	if !strings.HasPrefix(normalizedKey, "TERMINUS_") {
		normalizedKey = "TERMINUS_" + normalizedKey
	}

	// Handle specific configuration keys
	valueStr := fmt.Sprintf("%v", value)

	switch normalizedKey {
	case "TERMINUS_HOST":
		c.Host = valueStr
	case "TERMINUS_PORT":
		fmt.Sscanf(valueStr, "%d", &c.Port)
	case "TERMINUS_PROTOCOL":
		c.Protocol = valueStr
	case "TERMINUS_TIMEOUT":
		fmt.Sscanf(valueStr, "%d", &c.Timeout)
	case "TERMINUS_CACHE_DIR":
		c.CacheDir = c.expandPath(valueStr)
	case "TERMINUS_PLUGINS_DIR":
		c.PluginsDir = c.expandPath(valueStr)
	case "TERMINUS_TOKENS_DIR":
		c.TokensDir = c.expandPath(valueStr)
	case "TERMINUS_DATE_FORMAT":
		c.DateFormat = valueStr
	}

	// Store in values map
	c.values[normalizedKey] = value
}

// expandPath expands placeholders in paths
func (c *Config) expandPath(path string) string {
	// Replace common placeholders
	path = strings.ReplaceAll(path, "[[TERMINUS_USER_HOME]]", c.HomeDir)
	path = strings.ReplaceAll(path, "~", c.HomeDir)

	return path
}

// ensureDirectories creates necessary directories if they don't exist
func (c *Config) ensureDirectories() error {
	dirs := []string{
		c.CacheDir,
		c.TokensDir,
		c.PluginsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// Get returns a configuration value by key
func (c *Config) Get(key string) (interface{}, bool) {
	normalizedKey := strings.ToUpper(key)
	if !strings.HasPrefix(normalizedKey, "TERMINUS_") {
		normalizedKey = "TERMINUS_" + normalizedKey
	}

	value, ok := c.values[normalizedKey]
	return value, ok
}

// GetString returns a string configuration value
func (c *Config) GetString(key string) string {
	if value, ok := c.Get(key); ok {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

// GetInt returns an integer configuration value
func (c *Config) GetInt(key string) int {
	if value, ok := c.Get(key); ok {
		var result int
		fmt.Sscanf(fmt.Sprintf("%v", value), "%d", &result)
		return result
	}
	return 0
}

// GetBool returns a boolean configuration value
func (c *Config) GetBool(key string) bool {
	if value, ok := c.Get(key); ok {
		valueStr := strings.ToLower(fmt.Sprintf("%v", value))
		return valueStr == "true" || valueStr == "1" || valueStr == "yes"
	}
	return false
}

// GetBaseURL returns the full API base URL
func (c *Config) GetBaseURL() string {
	return fmt.Sprintf("%s://%s:%d/api", c.Protocol, c.Host, c.Port)
}
