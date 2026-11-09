package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvMapsDottedKeysAndPreservesShellOverrides(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, ".env")
	content := "postgresql.skeleton.host=dotenv-host\napps.name=dotenv-app\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write dotenv file: %v", err)
	}

	t.Setenv("PERSONALITIES_APPS_NAME", "shell-app")
	hostEnvName := "PERSONALITIES_POSTGRESQL_SKELETON_HOST"
	previousHost, hadPreviousHost := os.LookupEnv(hostEnvName)
	if err := os.Unsetenv(hostEnvName); err != nil {
		t.Fatalf("unset host environment: %v", err)
	}
	t.Cleanup(func() {
		if hadPreviousHost {
			_ = os.Setenv(hostEnvName, previousHost)
		} else {
			_ = os.Unsetenv(hostEnvName)
		}
	})

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("load dotenv: %v", err)
	}
	if got := os.Getenv(hostEnvName); got != "dotenv-host" {
		t.Fatalf("expected dotenv host to be mapped, got %q", got)
	}
	if got := os.Getenv("PERSONALITIES_APPS_NAME"); got != "shell-app" {
		t.Fatalf("expected shell override to be preserved, got %q", got)
	}
}
