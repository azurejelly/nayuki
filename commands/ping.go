package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c *PingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with 'Pong!' if the bot is online.",
	}
}

func (c *PingCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) (err error) {
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})

	return errors.New("test")
}
