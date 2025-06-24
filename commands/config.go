package commands

import (
	"fmt"

	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
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
					{
						Name:        "clear",
						Description: "Clears the selected suggestions channel for this server.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "logs",
				Description: "Configuration for the suggestions logs channel.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "get",
						Description: "Get the current suggestions logs channel for this server.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "set",
						Description: "Set the current suggestions logs channel for this server.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "channel",
								Description: "The suggestions logs channel to use for this server.",
								Type:        discordgo.ApplicationCommandOptionChannel,
								Required:    true,
							},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "threads",
				Description: "Enable or disable suggestion threads.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "enable",
						Description: "Enable thread creation for new suggestions.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "disable",
						Description: "Disable thread creation for new suggestions.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "status",
						Description: "Returns whether suggestion threads are currently enabled.",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
		},
	}
}

func (c *ConfigCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	opt := event.ApplicationCommandData().Options
	group := opt[0].Name
	action := opt[0].Options[0]

	utils.DeferEphemeral(s, i)
	server, err := database.GetOrCreateServer(i.GuildID)
	if err != nil {
		return utils.UpdateDeferredEphemeral(s, i, fmt.Sprintf(":x: Failed to fetch and/or create server data: ```\n%s\n```", err.Error()))
	}

	switch group {
	case "channel", "logs":
		switch action.Name {
		case "get":
			msg := ""

			if group == "logs" {
				if server.LogsChannel == "" {
					msg = ":x: This Discord server has no logging channel set."
				} else {
					msg = fmt.Sprintf(":information_source: The logging channel for this server is currently set to <#%s>", server.LogsChannel)
				}
			} else {
				if server.Channel == "" {
					msg = ":x: This Discord server has no suggestions channel set."
				} else {
					msg = fmt.Sprintf(":information_source: The suggestions channel for this server is currently set to <#%s>", server.Channel)
				}
			}

			return utils.UpdateDeferredEphemeral(s, i, msg)
		case "clear":
			msg := ""

			if group == "logs" {
				server.LogsChannel = ""
				msg = ":white_check_mark: Logging channel cleared."
			} else {
				server.Channel = ""
				msg = ":white_check_mark: Suggestions channel cleared."
			}

			err := database.SaveServer(server)
			if err != nil {
				msg = fmt.Sprintf(":x: Could not update server data: \n```\n%s\n```", err.Error())
			}

			return utils.UpdateDeferredEphemeral(s, i, msg)
		default:
			channel := action.GetOption("channel").ChannelValue(s)
			if channel.Type != discordgo.ChannelTypeGuildText {
				return utils.UpdateDeferredEphemeral(s, i, ":x: You must provide a text channel.")
			}

			if channel.GuildID != i.GuildID {
				return utils.UpdateDeferredEphemeral(s, i, ":x: Server ID mismatch.")
			}

			msg := ""

			if group == "logs" {
				server.LogsChannel = channel.ID
				msg = fmt.Sprintf(":white_check_mark: The logging channel for this server has been set to <#%s>", channel.ID)
			} else {
				server.Channel = channel.ID
				msg = fmt.Sprintf(":white_check_mark: The suggestions channel for this server has been set to <#%s>", channel.ID)
			}

			err = database.SaveServer(server)
			if err != nil {
				msg = fmt.Sprintf(":x: Failed to save new server data into database: ```\n%s\n```", err.Error())
			}

			return utils.UpdateDeferredEphemeral(s, i, msg)
		}
	case "threads":
		msg := ""

		switch action.Name {
		case "enable":
			server.CreateThreads = true
			msg = ":white_check_mark: Thread creation for new suggestions has been **enabled**."
		case "disable":
			msg = ":white_check_mark: Thread creation for new suggestions has been **disabled**."
			server.CreateThreads = false
		case "status":
			if server.CreateThreads {
				return utils.UpdateDeferredEphemeral(s, i, ":information_source: Thread creation for suggestions is currently **enabled**.")
			} else {
				return utils.UpdateDeferredEphemeral(s, i, ":information_source: Thread creation for suggestions is currently **disabled**.")
			}
		}

		err = database.SaveServer(server)
		if err != nil {
			msg = fmt.Sprintf(":x: Failed to save server data into database: ```\n%s\n```", err.Error())
			return utils.UpdateDeferredEphemeral(s, i, msg)
		}

		return utils.UpdateDeferredEphemeral(s, i, msg)
	}

	return nil
}
