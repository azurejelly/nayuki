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
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "response",
				Description: "An optional response for the person who submitted the suggestion.",
				Required:    false,
			},
		},
	}
}

func (c *DeclineCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	id := i.ApplicationCommandData().GetOption("id").StringValue()
	response := func() string {
		if opt := i.ApplicationCommandData().GetOption("response"); opt != nil {
			return opt.StringValue()
		}

		return ""
	}()

	return helper.TakeSuggestionAction(s, i, id, response, "declined", utils.NEGATIVE_EMBED_COLOR)
}
