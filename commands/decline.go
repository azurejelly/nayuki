package commands

import (
	"fmt"
	"log"

	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
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

// FIXME: duplicate code
func (c *DeclineCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	id := i.ApplicationCommandData().GetOption("id").StringValue()

	utils.DeferEphemeral(s, i)
	suggestion, err := database.FindSuggestion(id)

	if err != nil {
		return utils.UpdateDeferredEphemeral(s, i, fmt.Sprintf(":x: Failed to fetch suggestion by ID: \n```\n%s\n```", err.Error()))
	}

	if suggestion == nil {
		return utils.UpdateDeferredEphemeral(s, i, ":x: Could not find a suggestion with that ID.")
	}

	server, err := database.GetOrCreateServer(i.GuildID)
	if err != nil {
		return utils.UpdateDeferredEphemeral(s, i, fmt.Sprintf(":x: Failed to fetch and/or create server data: ```\n%s\n```", err.Error()))
	}

	err = database.DeleteSuggestion(suggestion)
	if err != nil {
		return utils.UpdateDeferredEphemeral(s, i, fmt.Sprintf(":x: Could not delete suggestion from database: ```\n%s\n```", err.Error()))
	}

	result := ":white_check_mark: Suggestion declined."
	msg, err := s.ChannelMessage(suggestion.Channel, suggestion.Message)
	likes, dislikes := 0, 0

	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was declined, but the original message could not be fetched:\n```\n%s\n```", err.Error())
	} else {
		for _, r := range msg.Reactions {
			if r.Emoji.Name == "👍" {
				likes += (r.Count - 1) // substract 1 since the bot always reacts once
			}

			if r.Emoji.Name == "👎" {
				dislikes += (r.Count - 1) // substract 1 since the bot always reacts once
			}
		}

		if msg.Thread != nil {
			_, err = s.ChannelEditComplex(msg.Thread.ID, &discordgo.ChannelEdit{
				Locked:   utils.Ptr(true),
				Archived: utils.Ptr(true),
			})

			if err != nil {
				result = fmt.Sprintf(":warning: The suggestion was declined, but its thread couldn't be locked/archived:\n```\n%s\n```", err.Error())
			}
		}
	}

	err = s.ChannelMessageDelete(suggestion.Channel, suggestion.Message)
	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was declined, but the original message couldn't be deleted: \n```\n%s\n```", err)
		log.Println("could not delete original suggestion message:", err)
	}

	// TODO: make the embed a bit better
	if server.LogsChannel != "" {
		title := suggestion.Title
		content := suggestion.Content

		if len(title) > 256 {
			title = title[:256]
		}

		if len(content) > 4096 {
			content = content[:4096]
		}

		embed := embed.NewEmbed()
		embed.SetTitle(title)
		embed.SetDescription(content)
		embed.AddField("Likes", fmt.Sprintf(":thumbs_up: %d", likes))
		embed.AddField("Dislikes", fmt.Sprintf(":thumbs_down: %d", dislikes))
		embed.AddField("Declined by", fmt.Sprintf("<@%s>", i.Member.User.ID))
		embed.AddField("Suggested by", fmt.Sprintf("<@%s> (%s)", suggestion.Author, suggestion.AuthorName))
		embed.SetFooter(fmt.Sprintf("ID: %s", suggestion.ID.Hex()))
		embed.SetColor(0xeb3731)
		embed.InlineAllFields()

		u, _ := s.User(suggestion.Author)
		if u != nil {
			embed.SetAuthor("Suggestion declined.", u.AvatarURL("1024"))
		} else {
			embed.SetAuthor("Suggestion declined.")
		}

		s.ChannelMessageSendEmbed(server.LogsChannel, embed.MessageEmbed)
	}

	return utils.UpdateDeferredEphemeral(s, i, result)
}
