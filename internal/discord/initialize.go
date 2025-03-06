package discord

import (
	"fmt"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal"

	"github.com/bwmarrin/discordgo"
)

// Initiates the Discord monitoring via the specified token
func Initialize() error {
	// Create a new Discord session
	dg, err := discordgo.New("Bot " + internal.Configuration.Discord.Token)
	if err != nil {
		return err
	}

	// Specify required intents for the session
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	dg.AddHandler(ready)        // Add a ready handler to notify once monitoring has started
	dg.AddHandler(startMonitor) // Add a message creation handler to reroute messages

	// Open a websocket connection
	if err := dg.Open(); err != nil {
		return err
	}
	defer dg.Close()

	// Wait for signals to exit
	select {}
}

// Notifies the user once monitoring has started
func ready(s *discordgo.Session, _ *discordgo.Ready) {
	fmt.Printf("Monitor running on %s... Press Ctrl-C to exit\n", s.State.User.Username)
}
