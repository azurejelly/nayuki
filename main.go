package main

import (
	"flag"

	"github.com/azurejelly/nayuki/config"
	"github.com/azurejelly/nayuki/database"
	"github.com/azurejelly/nayuki/discord"
)

var (
	token string
	uri   string
	db    string
)

func init() {
	flag.StringVar(&token, "t", "", "The Discord bot token to use for authentication.")
	flag.StringVar(&uri, "u", "", "The MongoDB URI the program should use. Used for storing suggestions.")
	flag.StringVar(&db, "d", "nayuki", "The database name. Defaults to 'nayuki'")
	flag.Parse()
}

func main() {
	config.Load()

	if token == "" {
		token = config.GetToken()
	}

	if uri == "" {
		uri = config.GetMongoURI()
	}

	if db == "" {
		db = config.GetMongoDatabase()
	}

	database.Init(uri, db)
	discord.Init(token)
}
