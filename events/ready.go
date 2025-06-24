package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("logged in as %s (id: %s)\n", s.State.User.Username, s.State.User.ID)
	s.UpdateWatchStatus(0, "ur cool suggestions")
}
