# Discord Mavely Router

A lightweight Discord bot that automatically detects links in messages from specified channels and replaces them with Mavely links. This bot supports configurable source and target channels, as well as an optional echo feature to resend modified messages in the same channel.

## Features
- **Link Detection & Replacement**: Scans messages in designated source channels and replaces all detected links with Mavely equivalents.
- **Flexible Channel Configuration**:
  - **Source Channels**: Monitor specific channels for messages containing links.
  - **Target Channels**: Send messages with replaced Mavely links to one or more designated channels.
  - **Echo Channels**: Optionally delete the original message and resend it in the same channel with Mavely link replacements.
- **Simple Setup**: Easy to configure and deploy for your Discord server.

## How It Works
1. The bot listens for messages in all specified **source channels**.
2. When a message containing a link is detected, it processes the message and generates a new version with Mavely links.
3. The modified message is then:
   - Sent to all **target channels**, and/or
   - Resent in the same channel (replacing the original message) if configured as an **echo channel**.

## Usage
- Define your source, target, and echo channels in the bot’s configuration file.
- Deploy the bot to your Discord server.
- Watch as it seamlessly converts links into Mavely links based on your setup!

## Configuration
The bot requires a JSON configuration file to operate. Below is an example configuration and the default config directory paths for each operating system.

### Example JSON Configuration
```json
{
  "guild_id": "123456789012345678",
  "echo_channels": ["channel_id_1", "channel_id_2"],
  "relay_source": ["source_channel_id_1", "source_channel_id_2"],
  "relay_target": ["target_channel_id_1", "target_channel_id_2"],
  "mavely": {
    "username": "your_mavely_username",
    "password": "your_mavely_password"
  },
  "discord": {
    "token": "your_discord_bot_token",
    "application_id": "your_discord_app_id"
  }
}
```

- **`guild_id`**: The ID of the Discord server (guild) where the bot operates.
- **`echo_channels`**: List of channel IDs where original messages are replaced with Mavely links.
- **`relay_source`**: List of channel IDs to monitor for links.
- **`relay_target`**: List of channel IDs where messages with Mavely links are sent.
- **`mavely`**: Credentials for authenticating with the Mavely service.
- **`discord`**: Discord bot token and application ID for authentication.

### Config Paths
After first run, the configuration file will be created in the appropriate user config directory for your operating system:
- **Linux**: `~/.config/discord-mavely-router/config.json`
- **macOS**: `~/Library/Application Support/discord-mavely-router/config.json`
- **Windows**: `%AppData%\discord-mavely-router\config.json` (e.g., `C:\Users\<YourUsername>\AppData\Roaming\...`)

## License
This project is licensed under the [MIT license](https://opensource.org/license/mit)
