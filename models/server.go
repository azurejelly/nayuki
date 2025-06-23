package models

import "github.com/kamva/mgm/v3"

type Server struct {
	mgm.DefaultModel `bson:",inline"`

	Guild   string `json:"guild_id" bson:"guild_id"`
	Channel string `json:"channel_id" bson:"channel_id"`
}

func NewServer(guild string, channel string) *Server {
	return &Server{
		Guild:   guild,
		Channel: channel,
	}
}
