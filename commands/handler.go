package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Command() *discordgo.ApplicationCommand
	Run(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

var Commands = []Command{
	&PingCommand{},
}
