package models

import "github.com/kamva/mgm/v3"

type Server struct {
	mgm.DefaultModel `bson:",inline"`

	Guild         string `json:"guild_id" bson:"guild_id"`
	Channel       string `json:"channel_id" bson:"channel_id"`
	LogsChannel   string `json:"logs_channel_id" bson:"logs_channel_id"`
	CreateThreads bool   `json:"create_threads" bson:"create_threads"`
}

func NewServer(guild string, channel string, logs string, threads bool) *Server {
	return &Server{
		Guild:         guild,
		Channel:       channel,
		LogsChannel:   logs,
		CreateThreads: threads,
	}
}

func DefaultServer(guild string) *Server {
	return NewServer(guild, "", "", false)
}
