package discord

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azurejelly/nayuki/commands"
	"github.com/azurejelly/nayuki/config"
	"github.com/azurejelly/nayuki/events"
	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session

func Init(token string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("failed to initialize discord session: ", err)
		return
	}

	session.AddHandler(events.InteractionCreate)
	session.AddHandler(events.Ready)

	if err := session.Open(); err != nil {
		log.Fatalln("failed to open discord session: ", err)
		return
	}

	defer session.Close()

	log.Printf("registering %d command(s)\n", len(commands.List))

	guild := config.GetGuildId()
	if guild != "" {
		log.Println("only registering commands for guild", guild)
	}

	for _, c := range commands.List {
		cmd := c.Command()

		_, err = session.ApplicationCommandCreate(session.State.User.ID, guild, cmd)
		if err != nil {
			log.Fatalf("failed to register command %s: %s", cmd.Name, err)
			return
		} else {
			log.Println("registered command", cmd.Name)
		}
	}

	log.Println("finished registering commands")

	// Run until we detect a stop signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("now exiting. goodbye.")
}

func GetSession() *discordgo.Session {
	return session
}
