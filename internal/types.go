package internal

type Config struct {
	GuildID      string   `json:"guild_id"`
	EchoChannels []string `json:"echo_channels"`
	RelaySource  []string `json:"relay_source"`
	RelayTarget  []string `json:"relay_target"`
	Mavely       Mavely   `json:"mavely"`
	Discord      Discord  `json:"discord"`
}

type Discord struct {
	Token         string `json:"token"`
	ApplicationID string `json:"application_id"`
}

type Mavely struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
