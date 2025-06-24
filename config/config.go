package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("failed to load .env file: ", err)
		return
	}
}

func GetGuildId() string {
	return os.Getenv("GUILD_ID")
}

func GetToken() string {
	return os.Getenv("TOKEN")
}

func GetMongoURI() string {
	return os.Getenv("MONGO_URI")
}

func GetMongoDatabase() string {
	db := os.Getenv("MONGO_DATABASE")
	if db == "" {
		return "nayuki"
	}

	return db
}
