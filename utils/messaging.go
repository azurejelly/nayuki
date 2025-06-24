package utils

import (
	"github.com/bwmarrin/discordgo"
)

func Reply(s *discordgo.Session, i *discordgo.Interaction, str string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
		},
	})
}

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
	})
}

func DeferEphemeral(s *discordgo.Session, i *discordgo.Interaction) error {
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
	})

	return err
}

func UpdateDeferredEphemeral(s *discordgo.Session, i *discordgo.Interaction, str string) error {
	_, err := s.FollowupMessageCreate(i, false, &discordgo.WebhookParams{
		Content: str,
		Flags:   discordgo.MessageFlagsEphemeral,
	})

	return err
}
