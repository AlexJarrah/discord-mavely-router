package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal"
)

// Returns the config directory for the OS
func getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, internal.APP_ID), nil
}

// Returns the data directory for the OS
func getDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", internal.APP_ID), nil
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", internal.APP_ID), nil
	case "linux":
		return filepath.Join(homeDir, ".local", "share", internal.APP_ID), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func GetConfigFile() (config internal.Config, err error) {
	path := filepath.Join(ConfigDirectory, "config.json")
	file, err := os.ReadFile(path)
	if err != nil {
		return internal.Config{}, err
	}

	if err = json.Unmarshal(file, &config); err != nil {
		return internal.Config{}, err
	}

	return config, nil
}

func WriteConfigFile(config internal.Config) (err error) {
	json_, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return os.WriteFile(filepath.Join(ConfigDirectory, "config.json"), json_, 0o677)
}
