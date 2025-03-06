package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Initialize() error {
	var err error
	if ConfigDirectory, err = getConfigDir(); err != nil {
		return err
	}

	if DataDirectory, err = getDataDir(); err != nil {
		return err
	}

	// Maps out the file's paths with their default data
	files := map[string][]byte{
		filepath.Join(ConfigDirectory, "config.json"): getDefaultConfigJSON(),
		filepath.Join(DataDirectory, "activity.log"):  nil,
		filepath.Join(DataDirectory, ".token"):        nil,
	}

	// Creates and writes the default data to the file if it does not exist
	for filename, data := range files {
		dir := filepath.Dir(filename)
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		if _, err = os.Stat(filename); os.IsNotExist(err) {
			if err = os.WriteFile(filename, data, 0o644); err != nil {
				return err
			}
		}
	}

	InitializeLogger()

	internal.Configuration, err = GetConfigFile()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// InitializeLogger sets up the log file rotation and formatter, printing to both file and stdout
func InitializeLogger() {
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(DataDirectory, "activity.log"),
		MaxSize:    100, // megabytes
		MaxBackups: 2,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Create a multi-writer that writes to both the file and stdout
	mw := io.MultiWriter(fileLogger, os.Stdout)

	log.SetOutput(mw)
	log.SetFlags(0)
}

func getDefaultConfigJSON() []byte {
	defaultConfig := internal.Config{
		GuildID: "Server ID",
		EchoChannels: []string{
			"Channels to mirror",
		},
		RelaySource: []string{
			"Channels to route to the relay target",
		},
		RelayTarget: []string{
			"Channels to send output of relay sources to",
		},
		Discord: internal.Discord{
			Token:         "Discord bot token",
			ApplicationID: "Discord application ID",
		},
		Mavely: internal.Mavely{
			Username: "Mavely username",
			Password: "Mavely password",
		},
	}

	res, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return []byte("{}")
	}

	return res
}
