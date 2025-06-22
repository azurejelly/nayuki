package commands

import (
	"github.com/azurejelly/nayuki/config"
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PingCommand struct{}

var (
	_ ken.SlashCommand       = (*PingCommand)(nil)
	_ ken.GuildScopedCommand = (*PingCommand)(nil)
)

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Version() string {
	return "1.0.0"
}

func (c *PingCommand) Description() string {
	return "Should return 'Pong!' if the bot is online."
}

func (c *PingCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *PingCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *PingCommand) Guild() string {
	return config.GetGuildId()
}

func (c *PingCommand) IsDmCapable() bool {
	return true
}

func (c *PingCommand) Run(ctx ken.Context) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})

	return
}
