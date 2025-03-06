package main

import (
	"log"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal/discord"
	"gitlab.com/AlexJarrah/discord-mavely-router/internal/filesystem"
)

func main() {
	if err := filesystem.Initialize(); err != nil {
		log.Fatalf("Failed to initialize necessary files: %v", err)
	}

	if err := discord.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Discord bot: %v", err)
	}
}
