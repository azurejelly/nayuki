package utils

import (
	"github.com/bwmarrin/discordgo"
)

func ReplyEmbed(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func ReplyEphemeral(s *discordgo.Session, i *discordgo.Interaction, str string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func Defer(s *discordgo.Session, i *discordgo.Interaction) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func UpdateDeferred(s *discordgo.Session, i *discordgo.Interaction, str string) error {
	_, err := s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: str,
		Flags:   discordgo.MessageFlagsEphemeral,
	})

	return err
}

func CountReactions(msg *discordgo.Message, name string, ignoreBot ...bool) int {
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

	if counter > 0 && len(ignoreBot) > 0 && ignoreBot[0] {
		counter--
	}

	return counter
}
