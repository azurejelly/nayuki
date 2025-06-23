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
	i := event.Interaction
	return utils.ReplyEphemeral(s, i, ":ping_pong: Pong!")
}
