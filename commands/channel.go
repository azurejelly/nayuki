package commands

import (
	"github.com/azurejelly/nayuki/models"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type ConfigCommand struct{}

func (c *ConfigCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config",
		Description: "Sets the suggestions channel that should be used in this server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "channel",
				Description: "Configuration for the preferred suggestions channel.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "get",
						Description: "Get the current suggestions channel for this server.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "set",
						Description: "Set the current suggestions channel for this server.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "channel",
								Description: "The suggestions channel to use for this server.",
								Type:        discordgo.ApplicationCommandOptionChannel,
								Required:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (c *ConfigCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	opt := event.ApplicationCommandData().Options

	switch opt[0].Name {
	case "channel":
		action := opt[0].Options[0]
		if action.Name == "get" {
			utils.DeferEphemeral(s, i)
			server := &models.Server{}
			coll := mgm.Coll(server)
			result := coll.FindOne(mgm.Ctx(), &bson.M{"guild_id": i.GuildID})
			result.Decode(&server)

			if server.Guild == "" || server.Channel == "" {
				return utils.UpdateDeferredEphemeral(s, i, ":no_entry: This Discord server has no suggestions channel set.")
			} else {
				msg := ":information_source: The suggestions channel for this server is currently set to <#" + server.Channel + ">"
				return utils.UpdateDeferredEphemeral(s, i, msg)
			}
		} else {
			channel := action.GetOption("channel").ChannelValue(s)
			if channel.Type != discordgo.ChannelTypeGuildText {
				return utils.ReplyEphemeral(s, i, ":no_entry: Only guild text channels can be used as a suggestions channel.")
			}

			utils.DeferEphemeral(s, i)
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

			msg := ":white_check_mark: The suggestions channel for this server has been set to <#" + channel.ID + ">"
			return utils.UpdateDeferredEphemeral(s, i, msg)
		}
	default:
		return utils.ReplyEphemeral(s, i, ":no_entry: Not implemented")
	}
}
