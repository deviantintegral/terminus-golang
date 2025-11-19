package commands

import (
	"testing"
)

func TestRedisEnableCmdStructure(t *testing.T) {
	// Test enable command structure
	if redisEnableCmd.Use != "redis:enable <site>" {
		t.Errorf("expected redisEnableCmd.Use to be 'redis:enable <site>', got '%s'", redisEnableCmd.Use)
	}

	if redisEnableCmd.Short == "" {
		t.Error("redisEnableCmd.Short should not be empty")
	}

	if redisEnableCmd.Long == "" {
		t.Error("redisEnableCmd.Long should not be empty")
	}
}

func TestRedisDisableCmdStructure(t *testing.T) {
	// Test disable command structure
	if redisDisableCmd.Use != "redis:disable <site>" {
		t.Errorf("expected redisDisableCmd.Use to be 'redis:disable <site>', got '%s'", redisDisableCmd.Use)
	}

	if redisDisableCmd.Short == "" {
		t.Error("redisDisableCmd.Short should not be empty")
	}

	if redisDisableCmd.Long == "" {
		t.Error("redisDisableCmd.Long should not be empty")
	}
}

func TestRedisCommands(t *testing.T) {
	expectedCommands := []string{"redis:enable", "redis:disable"}

	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == expected || (len(cmd.Use) > len(expected) && cmd.Use[:len(expected)] == expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected command '%s' not found in rootCmd", expected)
		}
	}
}

func TestRedisEnableRequiresAuth(t *testing.T) {
	// Save old context and create a minimal context without session
	oldContext := cliContext
	cliContext = &CLIContext{
		APIClient: nil,
	}
	defer func() { cliContext = oldContext }()

	err := runRedisEnable(nil, []string{"test-site"})
	if err == nil {
		t.Error("expected error when not authenticated")
	}
}

func TestRedisDisableRequiresAuth(t *testing.T) {
	// Save old context and create a minimal context without session
	oldContext := cliContext
	cliContext = &CLIContext{
		APIClient: nil,
	}
	defer func() { cliContext = oldContext }()

	err := runRedisDisable(nil, []string{"test-site"})
	if err == nil {
		t.Error("expected error when not authenticated")
	}
}

func TestRedisEnableRequiresArg(t *testing.T) {
	// Test that enable command requires exactly one argument
	cmd := redisEnableCmd

	// The command should have Args set to ExactArgs(1)
	if cmd.Args == nil {
		t.Error("redisEnableCmd should have Args validator set")
	}
}

func TestRedisDisableRequiresArg(t *testing.T) {
	// Test that disable command requires exactly one argument
	cmd := redisDisableCmd

	// The command should have Args set to ExactArgs(1)
	if cmd.Args == nil {
		t.Error("redisDisableCmd should have Args validator set")
	}
}
