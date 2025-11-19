package commands

import (
	"testing"
)

func TestRedisCmdStructure(t *testing.T) {
	// Test command structure
	if redisCmd.Use != "redis" {
		t.Errorf("expected redisCmd.Use to be 'redis', got '%s'", redisCmd.Use)
	}

	if redisCmd.Short == "" {
		t.Error("redisCmd.Short should not be empty")
	}

	if redisCmd.Long == "" {
		t.Error("redisCmd.Long should not be empty")
	}
}

func TestRedisEnableCmdStructure(t *testing.T) {
	// Test enable command structure
	if redisEnableCmd.Use != "enable <site>" {
		t.Errorf("expected redisEnableCmd.Use to be 'enable <site>', got '%s'", redisEnableCmd.Use)
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
	if redisDisableCmd.Use != "disable <site>" {
		t.Errorf("expected redisDisableCmd.Use to be 'disable <site>', got '%s'", redisDisableCmd.Use)
	}

	if redisDisableCmd.Short == "" {
		t.Error("redisDisableCmd.Short should not be empty")
	}

	if redisDisableCmd.Long == "" {
		t.Error("redisDisableCmd.Long should not be empty")
	}
}

func TestRedisSubcommands(t *testing.T) {
	// Test that enable and disable are subcommands of redis
	subcommands := redisCmd.Commands()

	expectedCommands := map[string]bool{
		"enable":  false,
		"disable": false,
	}

	for _, cmd := range subcommands {
		if _, exists := expectedCommands[cmd.Name()]; exists {
			expectedCommands[cmd.Name()] = true
		}
	}

	for cmdName, found := range expectedCommands {
		if !found {
			t.Errorf("expected '%s' to be a subcommand of redis", cmdName)
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

func TestRedisCommandsRegistered(t *testing.T) {
	// Verify redis commands are properly registered
	// Find redis in root commands
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() != "redis" {
			continue
		}

		found = true

		// Check subcommands
		subCmds := cmd.Commands()
		hasEnable := false
		hasDisable := false

		for _, subCmd := range subCmds {
			switch subCmd.Name() {
			case "enable":
				hasEnable = true
			case "disable":
				hasDisable = true
			}
		}

		if !hasEnable {
			t.Error("redis command should have 'enable' subcommand")
		}
		if !hasDisable {
			t.Error("redis command should have 'disable' subcommand")
		}

		break
	}

	if !found {
		t.Error("redis command should be registered with root command")
	}
}
