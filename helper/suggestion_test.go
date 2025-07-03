package helper

import (
	"testing"

	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
)

func TestCountReactions(t *testing.T) {
	msg := &discordgo.Message{
		Reactions: []*discordgo.MessageReactions{
			{
				Emoji: &discordgo.Emoji{
					Name: utils.LIKE_EMOJI,
				},
				Count: 5,
			},
			{
				Emoji: &discordgo.Emoji{
					Name: utils.DISLIKE_EMOJI,
				},
				Count: 1,
			},
			{
				Emoji: &discordgo.Emoji{
					Name: "example",
				},
				Count: 0,
			},
			{
				Emoji: &discordgo.Emoji{
					Name: "another",
				},
				Count: -9,
			},
		},
	}

	if got := countReactions(msg, utils.LIKE_EMOJI); got != 4 {
		t.Errorf("countReactions(msg, %s) = %d; expected 4 instead\n", utils.LIKE_EMOJI, got)
	}

	if got := countReactions(msg, utils.DISLIKE_EMOJI); got != 0 {
		t.Errorf("countReactions(msg, %s) = %d; expected 0 instead\n", utils.DISLIKE_EMOJI, got)
	}

	if got := countReactions(msg, "example"); got != 0 {
		t.Errorf("countReactions(msg, %s) = %d; expected 0 instead\n", "example", got)
	}

	if got := countReactions(msg, "another"); got != 0 {
		t.Errorf("countReactions(msg, %s) = %d; expected 0 instead\n", "another", got)
	}

	if got := countReactions(nil, utils.DISLIKE_EMOJI); got != 0 {
		t.Errorf("countReactions(nil, %s) = %d; expected 0 instead\n", utils.DISLIKE_EMOJI, got)
	}
}
