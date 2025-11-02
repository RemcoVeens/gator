package config_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RemcoVeens/gator/internal/config"
	"github.com/stretchr/testify/assert"
)

func setupTestEnv(t *testing.T) string {
	tempDir := t.TempDir()

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
	return filepath.Join(tempDir, ".gatorconfig.json")
}

func Test_getConfigFilePath(t *testing.T) {
	expectedPath := setupTestEnv(t)
	cfg := config.Config{DBUrl: "test_url", CurentUserName: "test_user"}
	err := cfg.SetUser("test_user")
	assert.NoError(t, err, "SetUser should not return an error")

	_, err = os.Stat(expectedPath)
	assert.False(t, os.IsNotExist(err), "Config file should exist at the expected path")
}

func TestRead_Success(t *testing.T) {
	configPath := setupTestEnv(t)

	expectedConfig := config.Config{
		DBUrl:          "postgres://user:pass@host:5432/db",
		CurentUserName: "alice_doe",
	}
	validJSON := fmt.Sprintf(`{"db_url": "%s", "current_user_name": "%s"}`,
		expectedConfig.DBUrl, expectedConfig.CurentUserName)

	err := os.WriteFile(configPath, []byte(validJSON), 0644)
	assert.NoError(t, err, "Failed to write initial config file")

	actualConfig, err := config.Read()

	assert.NoError(t, err, "Read should not return an error for a valid file")
	assert.Equal(t, expectedConfig.DBUrl, actualConfig.DBUrl)
	assert.Equal(t, expectedConfig.CurentUserName, actualConfig.CurentUserName)
}

func TestRead_FileDoesNotExist(t *testing.T) {
	setupTestEnv(t)
	_, err := config.Read()

	assert.Error(t, err, "Read should return an error when file does not exist")
	assert.Contains(t, err.Error(), "failed to open config file", "Error message should indicate file opening failure")
	assert.True(t, errors.Is(err, os.ErrNotExist), "Error should wrap os.ErrNotExist")
}

func TestRead_InvalidJSON(t *testing.T) {
	configPath := setupTestEnv(t)
	invalidJSON := `{"db_url": "test_url", "current_user_name": "test_user"`
	err := os.WriteFile(configPath, []byte(invalidJSON), 0644)

	assert.NoError(t, err, "Failed to write invalid config file")

	_, err = config.Read()

	assert.Error(t, err, "Read should return an error for invalid JSON")
	assert.Contains(t, err.Error(), "failed to unmarshal JSON", "Error message should indicate JSON unmarshal failure")
}

func TestSetUser_Success(t *testing.T) {
	configPath := setupTestEnv(t)

	cfgToWrite := config.Config{
		DBUrl:          "sqlite:///test.db",
		CurentUserName: "bob_smith",
	}

	err := cfgToWrite.SetUser("bob_smith")
	assert.NoError(t, err, "SetUser should not return an error")

	byteValue, err := os.ReadFile(configPath)
	assert.NoError(t, err, "Should be able to read the file written by SetUser")
	var actualConfig config.Config
	err = json.Unmarshal(byteValue, &actualConfig)

	assert.NoError(t, err, "Content written by SetUser should be valid JSON")
	assert.Equal(t, cfgToWrite.DBUrl, actualConfig.DBUrl)
	assert.Equal(t, cfgToWrite.CurentUserName, actualConfig.CurentUserName)
}

func TestSetUser_CheckFilePermissions(t *testing.T) {
	configPath := setupTestEnv(t)
	cfg := config.Config{CurentUserName: "permission_check"}
	err := cfg.SetUser("bob_smith")
	assert.NoError(t, err)

	info, err := os.Stat(configPath)
	assert.NoError(t, err)

	expectedPerms := os.FileMode(0644)
	actualPerms := info.Mode().Perm()

	assert.Equal(t, expectedPerms, actualPerms, "File permissions should be 0644")
}
