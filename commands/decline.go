package commands

import (
	"github.com/azurejelly/nayuki/helper"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
)

type DeclineCommand struct{}

func (c *DeclineCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     "decline",
		Description:              "Declines a suggestion.",
		DefaultMemberPermissions: utils.Ptr(int64(discordgo.PermissionManageMessages)),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the suggestion.",
				Required:    true,
			},
		},
	}
}

func (c *DeclineCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	id := i.ApplicationCommandData().GetOption("id").StringValue()
	return helper.TakeSuggestionAction(s, i, id, "declined", utils.POSITIVE_EMBED_COLOR)
}
