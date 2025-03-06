package discord

import (
	"log"
	"slices"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal"

	"github.com/bwmarrin/discordgo"
)

func startMonitor(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID != internal.Configuration.GuildID || m.Author.ID == internal.Configuration.Discord.ApplicationID {
		return
	}

	if slices.Contains(internal.Configuration.EchoChannels, m.ChannelID) {
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			log.Println("Failed to delete message:", err)
		}
		if err := sendMessage(s, m.Message, m.ChannelID); err != nil {
			log.Println("Failed to send message:", err)
		}
	} else if slices.Contains(internal.Configuration.RelaySource, m.ChannelID) {
		for _, c := range internal.Configuration.RelayTarget {
			if err := sendMessage(s, m.Message, c); err != nil {
				log.Println("Failed to send message:", err)
			}
		}
	}
}
