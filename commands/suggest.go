package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/models"
	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/kamva/mgm/v3"
)

type SuggestCommand struct{}

func (c *SuggestCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "suggest",
		Description: "Make a suggestion to the administrators of this Discord server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "The title for this suggestion.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "description",
				Description: "The description for this suggestion.",
				Required:    true,
			},
		},
	}
}

func (c *SuggestCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) (err error) {
	i := event.Interaction

	utils.Defer(s, i)
	server, err := database.GetServer(i.GuildID)
	if err != nil {
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Failed to fetch and/or create server data: ```\n%s\n```", err.Error()))
	}

	if server == nil || server.Channel == "" {
		return utils.UpdateDeferred(s, i, ":x: Sorry, this server isn't accepting suggestions at the moment.")
	}

	channel := server.Channel
	msg, err := s.ChannelMessageSend(channel, "A new suggestion is just around the corner!")
	if err != nil {
		return utils.UpdateDeferred(s, i, ":x: Sorry, this server isn't accepting suggestions at the moment.")
	}

	data := i.ApplicationCommandData()
	title := utils.Truncate(data.GetOption("title").StringValue(), utils.MAX_TITLE_LENGTH)
	content := utils.Truncate(data.GetOption("description").StringValue(), utils.MAX_DESCRIPTION_LENGTH)

	suggestion := models.NewSuggestion(event.Member.User.ID, event.Member.User.Username, title, content, channel, msg.ID)
	coll := mgm.Coll(suggestion)
	err = coll.Create(suggestion)

	if err != nil {
		s.ChannelMessageDelete(channel, suggestion.Message)
		return utils.UpdateDeferred(s, i, fmt.Sprintf(":x: Failed to create suggestion: ```\n%s\n```", err.Error()))
	}

	embed := embed.NewEmbed()
	embed.SetTitle(title)
	embed.SetDescription(content)
	embed.SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("ID: %s", suggestion.ID.Hex()))
	embed.SetColor(0xe28de3)
	embed.Timestamp = time.Now().Format(time.RFC3339)

	s.ChannelMessageEditEmbed(channel, msg.ID, embed.MessageEmbed)
	s.ChannelMessageEdit(channel, msg.ID, fmt.Sprintf("New suggestion from <@%s>:", i.Member.User.ID))
	s.MessageReactionAdd(channel, msg.ID, "👍") // thumbs up
	s.MessageReactionAdd(channel, msg.ID, "👎") // thumbs down

	if server.CreateThreads {
		thread := fmt.Sprintf("%s - %s", suggestion.ID.Hex(), title)
		if len(thread) > 100 {
			thread = thread[:100]
		}

		_, err = s.MessageThreadStart(channel, msg.ID, thread, 1440 /* 1d */)
		if err != nil {
			log.Println("failed to create thread: ", err)
		}
	}

	return utils.UpdateDeferred(s, i, fmt.Sprintf(":white_check_mark: Suggestion created! The ID for it is `%s`.", suggestion.ID.Hex()))
}
