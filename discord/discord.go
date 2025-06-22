package discord

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azurejelly/nayuki/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/store"
)

var session *discordgo.Session
var k *ken.Ken

func Init(token string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("failed to initialize discord session: ", err)
		return
	}

	if err := session.Open(); err != nil {
		log.Fatal("failed to open discord session: ", err)
		return
	}

	defer session.Close()

	k, err := ken.New(session, ken.Options{
		CommandStore: store.NewDefault(),
	})

	if err != nil {
		log.Fatal("failed to create ken instance: ", err)
		return
	}

	k.RegisterCommands(new(commands.PingCommand))

	defer k.Unregister()

	log.Println("discord connection was successful")

	// Run until we detect a stop signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("now exiting. goodbye.")
}

func GetSession() *discordgo.Session {
	return session
}

func GetKen() *ken.Ken {
	return k
}
