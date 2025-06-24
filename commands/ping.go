package commands

import (
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c *PingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with 'Pong!' if the bot is online.",
	}
}

func (c *PingCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) (err error) {
	return utils.ReplyEphemeral(s, event.Interaction, ":ping_pong: Pong!")
}
