package models

import "github.com/kamva/mgm/v3"

type Suggestion struct {
	mgm.DefaultModel `bson:",inline"`

	Author      string       `json:"author_id" bson:"author_id"`
	Title       string       `json:"title" bson:"title"`
	Content     string       `json:"content" bson:"content"`
	Channel     string       `json:"channel_id" bson:"channel_id"`
	Message     string       `json:"message_id" bson:"message_id"`
	Attachments []Attachment `json:"attachments" bson:"attachments"`
}

func NewSuggestion(author string, title string, content string, channel string, message string) *Suggestion {
	return &Suggestion{
		Author:  author,
		Title:   title,
		Content: content,
		Channel: channel,
		Message: message,
	}
}
