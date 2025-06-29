package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Run(s *discordgo.Session, event *discordgo.InteractionCreate) error
}

var List = []Command{
	&PingCommand{},
	&ConfigCommand{},
	&SuggestCommand{},
	&AcceptCommand{},
	&DeclineCommand{},
	&InfoCommand{},
}
