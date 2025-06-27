package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type AcceptCommand struct{}

func (c *AcceptCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     "accept",
		Description:              "Accepts a suggestion.",
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
func (c *AcceptCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction
	id := i.ApplicationCommandData().GetOption("id").StringValue()
	utils.Defer(s, i)

	suggestion, err := database.FindSuggestion(id)
	if err != nil || suggestion == nil {
		return utils.UpdateDeferred(s, i, ":x: Could not find a suggestion with that ID.")
	}

	server, err := database.GetOrCreateServer(i.GuildID)
	if err != nil {
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Failed to fetch and/or create server data: ```\n%s\n```", err.Error()))
	}

	err = database.DeleteSuggestion(suggestion)
	if err != nil {
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Could not delete suggestion from database: ```\n%s\n```", err.Error()))
	}

	result := ":white_check_mark: Suggestion accepted!"
	msg, err := s.ChannelMessage(suggestion.Channel, suggestion.Message)
	likes, dislikes := 0, 0

	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was accepted, but the original message could not be fetched:\n```\n%s\n```", err.Error())
	} else {
		// Count likes and dislikes in the original message
		likes = utils.CountReactions(msg, "👍", true)
		dislikes = utils.CountReactions(msg, "👎", true)

		// Lock the original discussion thread, if available
		if msg.Thread != nil {
			_, err := s.ChannelEditComplex(msg.Thread.ID, &discordgo.ChannelEdit{
				Locked:   utils.Ptr(true),
				Archived: utils.Ptr(true),
			})

			if err != nil {
				result = fmt.Sprintf(":warning: The suggestion was accepted, but its thread couldn't be locked/archived:\n```\n%s\n```", err.Error())
			}
		}
	}

	// Delete the original message
	err = s.ChannelMessageDelete(suggestion.Channel, suggestion.Message)
	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was accepted, but the original message couldn't be deleted: \n```\n%s\n```", err)
		log.Println("could not delete original suggestion message:", err)
	}

	// TODO: make the embed a bit better
	if server.LogsChannel != "" {
		title := utils.Truncate(suggestion.Title, utils.MAX_TITLE_LENGTH)
		content := utils.Truncate(suggestion.Content, utils.MAX_DESCRIPTION_LENGTH)

		embed := embed.NewEmbed()
		embed.SetTitle(title)
		embed.SetDescription(content)
		embed.AddField("Likes", fmt.Sprintf(":thumbs_up: %d", likes))
		embed.AddField("Dislikes", fmt.Sprintf(":thumbs_down: %d", dislikes))
		embed.AddField("Approved by", fmt.Sprintf("<@%s>", i.Member.User.ID))
		embed.AddField("Suggested by", fmt.Sprintf("<@%s> (%s)", suggestion.Author, suggestion.AuthorName))
		embed.SetFooter(fmt.Sprintf("ID: %s", suggestion.ID.Hex()))
		embed.SetColor(0x0fd15c)
		embed.InlineAllFields()

		u, _ := s.User(suggestion.Author)
		if u != nil {
			embed.SetAuthor("Suggestion accepted!", u.AvatarURL("1024"))
		} else {
			embed.SetAuthor("Suggestion accepted!")
		}

		embed.Timestamp = time.Now().Format(time.RFC3339)
		s.ChannelMessageSendEmbed(server.LogsChannel, embed.MessageEmbed)
	}

	return utils.UpdateDeferred(s, i, result)
}
