package commands

import (
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
)

type SuggestCommand struct{}

func (c *SuggestCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "suggest",
		Description: "Make a suggestion to the administrators of this Discord server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "title",
				Description: "The title for this suggestion.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "description",
				Description: "The description for this suggestion.",
				Required:    true,
			},
		},
	}
}

func (c *SuggestCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) (err error) {
	i := event.Interaction
	return utils.ReplyEphemeral(s, i, ":no_entry: Not implemented")
}
