package discord

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gitlab.com/AlexJarrah/discord-mavely-router/internal/mavely"

	"github.com/bwmarrin/discordgo"
)

// replaceLinksInText replaces all links in the given text with Mavely links using the provided token.
func replaceLinksInText(text, token string) string {
	// Replace markdown links: [text](url)
	reMarkdown := regexp.MustCompile(`\[([^\]]*)\]\((\S+)\)`)
	text = reMarkdown.ReplaceAllStringFunc(text, func(match string) string {
		submatches := reMarkdown.FindStringSubmatch(match)
		if len(submatches) == 3 {
			textPart := submatches[1]
			url := submatches[2]
			mavelyURL, _, err := mavely.CreateLink(token, url)
			if err != nil {
				log.Printf("Failed to create Mavely link for %s: %v", url, err)
				// If Mavely link creation fails, keep the original markdown link
				return match
			}
			return "[" + textPart + "](" + mavelyURL + ")"
		}
		return match
	})

	// Replace plain URLs: http://... or https://...
	reURL := regexp.MustCompile(`https?://\S+`)
	text = reURL.ReplaceAllStringFunc(text, func(url string) string {
		mavelyURL, _, err := mavely.CreateLink(token, url)
		if err != nil {
			log.Printf("Failed to create Mavely link for %s: %v", url, err)
			// If Mavely link creation fails, keep the original URL
			return url
		}
		return mavelyURL
	})

	return text
}

// processEmbed creates a new embed with all links replaced by Mavely links using the provided token.
func processEmbed(embed *discordgo.MessageEmbed, token string) *discordgo.MessageEmbed {
	newEmbed := &discordgo.MessageEmbed{
		URL:         replaceURL(embed.URL, token),
		Title:       replaceLinksInText(embed.Title, token),
		Description: replaceLinksInText(embed.Description, token),
		Type:        embed.Type,
		Color:       embed.Color,
		Timestamp:   embed.Timestamp,
	}

	// Process Footer
	if embed.Footer != nil {
		newEmbed.Footer = &discordgo.MessageEmbedFooter{
			Text:         replaceLinksInText(embed.Footer.Text, token),
			IconURL:      replaceURL(embed.Footer.IconURL, token),
			ProxyIconURL: embed.Footer.ProxyIconURL,
		}
	}

	// Process Image
	if embed.Image != nil {
		newEmbed.Image = &discordgo.MessageEmbedImage{
			URL:      replaceURL(embed.Image.URL, token),
			ProxyURL: embed.Image.ProxyURL,
			Width:    embed.Image.Width,
			Height:   embed.Image.Height,
		}
	}

	// Process Thumbnail
	if embed.Thumbnail != nil {
		newEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL:      replaceURL(embed.Thumbnail.URL, token),
			ProxyURL: embed.Thumbnail.ProxyURL,
			Width:    embed.Thumbnail.Width,
			Height:   embed.Thumbnail.Height,
		}
	}

	// Process Video
	if embed.Video != nil {
		newEmbed.Video = &discordgo.MessageEmbedVideo{
			URL:    replaceURL(embed.Video.URL, token),
			Width:  embed.Video.Width,
			Height: embed.Video.Height,
		}
	}

	// Process Provider
	if embed.Provider != nil {
		newEmbed.Provider = &discordgo.MessageEmbedProvider{
			Name: embed.Provider.Name,
			URL:  replaceURL(embed.Provider.URL, token),
		}
	}

	// Process Author
	if embed.Author != nil {
		newEmbed.Author = &discordgo.MessageEmbedAuthor{
			Name:    embed.Author.Name,
			URL:     replaceURL(embed.Author.URL, token),
			IconURL: replaceURL(embed.Author.IconURL, token),
		}
	}

	// Process Fields
	for _, field := range embed.Fields {
		newField := &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  replaceLinksInText(field.Value, token),
			Inline: field.Inline,
		}
		newEmbed.Fields = append(newEmbed.Fields, newField)
	}

	return newEmbed
}

// processComponents processes message components to replace link button URLs with Mavely links.
func processComponents(components []discordgo.MessageComponent, token string) []discordgo.MessageComponent {
	var newComponents []discordgo.MessageComponent
	for _, comp := range components {
		switch c := comp.(type) {
		case *discordgo.ActionsRow:
			newRow := &discordgo.ActionsRow{}
			for _, subComp := range c.Components {
				newSubComp := processComponent(subComp, token)
				newRow.Components = append(newRow.Components, newSubComp)
			}
			newComponents = append(newComponents, newRow)
		default:
			newComponents = append(newComponents, comp)
		}
	}
	return newComponents
}

// processComponent processes individual components, replacing URLs in link buttons.
func processComponent(comp discordgo.MessageComponent, token string) discordgo.MessageComponent {
	switch c := comp.(type) {
	case *discordgo.Button:
		if c.Style == discordgo.LinkButton {
			newButton := *c
			newButton.URL = replaceURL(c.URL, token)
			return &newButton
		}
		return c
	default:
		return c
	}
}

// replaceURL replaces a single URL with a Mavely link using the provided token.
func replaceURL(url, token string) string {
	if url == "" {
		return url
	}
	mavelyURL, _, err := mavely.CreateLink(token, url)
	if err != nil {
		log.Printf("Failed to create Mavely link for %s: %v", url, err)
		// If Mavely link creation fails, return the original URL
		return url
	}
	return mavelyURL
}

// sendMessage sends a message with all links replaced by Mavely links using the provided token.
func sendMessage(s *discordgo.Session, msg *discordgo.Message, channel string) error {
	tokenData, err := mavely.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get Mavely token: %w", err)
	}
	token := tokenData.AccessToken

	// Process the message content
	newContent := replaceLinksInText(msg.Content, token)

	// Process all embeds
	var newEmbeds []*discordgo.MessageEmbed
	for _, embed := range msg.Embeds {
		newEmbed := processEmbed(embed, token)
		newEmbeds = append(newEmbeds, newEmbed)
	}

	// Process all components
	newComponents := processComponents(msg.Components, token)

	// Handle attachments as files (unchanged from original message)
	files := func() (res []*discordgo.File) {
		for _, f := range msg.Attachments {
			res = append(res, &discordgo.File{
				Name:        f.Filename,
				ContentType: f.ContentType,
				Reader:      strings.NewReader(msg.Content),
			})
		}
		return res
	}()

	// Send the modified message
	_, err = s.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
		Content:    newContent,
		Embeds:     newEmbeds,
		Components: newComponents,
		Files:      files,
	})
	return err
}
