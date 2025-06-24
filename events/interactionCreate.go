package events

import (
	"log"

	"github.com/azurejelly/nayuki/commands"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	for _, c := range commands.Commands {
		name := c.Command().Name
		if name == i.ApplicationCommandData().Name {
			log.Printf("user '%s' has executed '%s'", i.Member.User.Username, name)
			err := c.Run(s, i)
			if err != nil {
				log.Printf("an error occurred while running %s: %s", name, err)
			}
		}
	}
}
