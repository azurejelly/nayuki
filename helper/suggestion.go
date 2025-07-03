package helper

import (
	"fmt"
	"log"
	"time"

	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

func TakeSuggestionAction(
	s *discordgo.Session,
	i *discordgo.Interaction,
	id string,
	verb string,
	color int,
) error {
	utils.Defer(s, i)

	suggestion, err := database.FindSuggestion(id)
	if err != nil || suggestion == nil {
		return utils.UpdateDeferred(s, i, ":x: Could not find a suggestion with that ID.")
	}

	server, err := database.GetOrCreateServer(i.GuildID)
	if err != nil {
		log.Println("something went wrong while communicating with the database:", err)
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Failed to fetch and/or create server data: ```\n%s\n```", err))
	}

	err = database.DeleteSuggestion(suggestion)
	if err != nil {
		log.Println("something went wrong while communicating with the database:", err)
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Could not delete suggestion from database: ```\n%s\n```", err))
	}

	result := fmt.Sprintf(":white_check_mark: Suggestion %s!", verb)
	msg, err := s.ChannelMessage(suggestion.Channel, suggestion.Message)
	likes, dislikes := 0, 0

	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was %s, but the original message could not be fetched:\n```\n%s\n```", verb, err)
	} else {
		// Count likes and dislikes in the original message
		likes = countReactions(msg, utils.LIKE_EMOJI)
		dislikes = countReactions(msg, utils.DISLIKE_EMOJI)

		// Lock the original discussion thread, if available
		if msg.Thread != nil {
			_, err := s.ChannelEditComplex(msg.Thread.ID, &discordgo.ChannelEdit{
				Locked:   utils.Ptr(true),
				Archived: utils.Ptr(true),
			})

			if err != nil {
				result = fmt.Sprintf(":warning: The suggestion was %s, but its thread couldn't be locked/archived:\n```\n%s\n```", verb, err)
			}
		}
	}

	// Delete the original message
	err = s.ChannelMessageDelete(suggestion.Channel, suggestion.Message)
	if err != nil {
		result = fmt.Sprintf(":warning: The suggestion was %s, but the original message couldn't be deleted: \n```\n%s\n```", verb, err)
		log.Println("could not delete original suggestion message:", err)
	}

	// TODO: make the embed a bit better
	if server.LogsChannel != "" {
		title := utils.Truncate(suggestion.Title, utils.MAX_TITLE_LENGTH)
		content := utils.Truncate(suggestion.Content, utils.MAX_DESCRIPTION_LENGTH)

		embed := embed.NewEmbed()
		embed.SetAuthor(fmt.Sprintf("Suggestion %s!", verb))
		embed.SetTitle(title)
		embed.SetDescription(content)
		embed.AddField("Likes", fmt.Sprintf(":thumbs_up: %d", likes))
		embed.AddField("Dislikes", fmt.Sprintf(":thumbs_down: %d", dislikes))
		embed.AddField("Reviewer", fmt.Sprintf("<@%s>", i.Member.User.ID))
		embed.AddField("Suggested by", fmt.Sprintf("<@%s> (%s)", suggestion.Author, suggestion.AuthorName))
		embed.SetFooter(fmt.Sprintf("ID: %s", suggestion.ID.Hex()))
		embed.SetColor(color)
		embed.InlineAllFields()

		u, _ := s.User(suggestion.Author)
		if u != nil {
			embed.SetAuthor(embed.Author.Name, u.AvatarURL("1024"))
		}

		embed.Timestamp = time.Now().Format(time.RFC3339)
		s.ChannelMessageSendEmbed(server.LogsChannel, embed.MessageEmbed)
	}

	return utils.UpdateDeferred(s, i, result)
}

func countReactions(msg *discordgo.Message, name string) int {
	counter := 0

	if msg == nil {
		return counter
	}

	for _, r := range msg.Reactions {
		if r.Emoji.Name == name {
			counter += r.Count
			break
		}
	}

	if counter > 0 {
		counter--
	}

	if counter < 0 {
		counter = 0
	}

	return counter
}
