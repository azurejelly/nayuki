package commands

import "github.com/bwmarrin/discordgo"

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

func (c *SuggestCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":no_entry: Not implemented!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
