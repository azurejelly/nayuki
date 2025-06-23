package commands

import (
	"github.com/azurejelly/nayuki/models"
	"github.com/bwmarrin/discordgo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type SetChannelCommand struct{}

func (c *SetChannelCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "set-channel",
		Description: "Sets the suggestions channel that should be used in this server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The suggestions channel that should be used in this server.",
				Required:    true,
			},
		},
	}
}

func (c *SetChannelCommand) Run(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	interaction := i.Interaction
	channel := i.ApplicationCommandData().GetOption("channel").ChannelValue(s)

	if channel.Type != discordgo.ChannelTypeGuildText {
		return s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":no_entry: Only guild text channels can be used as a suggestions channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	server := &models.Server{}
	coll := mgm.Coll(server)
	result := coll.FindOne(mgm.Ctx(), &bson.M{"guild_id": i.GuildID})
	result.Decode(&server)

	// TODO: find a better way (if any) of checking if this guild has no server data at all
	if server.Guild == "" {
		server = models.NewServer(i.GuildID, channel.ID)
		coll.Create(server)
	} else {
		server.Channel = channel.ID
		coll.Update(server)
	}

	_, err := s.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: ":white_check_mark: The suggestions channel for this server has been set to <#" + channel.ID + ">",
	})

	return err
}
