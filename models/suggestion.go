package models

import "github.com/kamva/mgm/v3"

type Suggestion struct {
	mgm.DefaultModel `bson:",inline"`

	Author      string       `json:"author_id" bson:"author_id"`
	Title       string       `json:"title" bson:"title"`
	Content     string       `json:"content" bson:"content"`
	Message     string       `json:"message_id" bson:"message_id"`
	Attachments []Attachment `json:"attachments" bson:"attachments"`
}

type Attachment struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

func NewSuggestion(author string, title string, content string, message string) *Suggestion {
	return &Suggestion{
		Author:  author,
		Title:   title,
		Content: content,
		Message: message,
	}
}

func NewAttachment(name string, url string) *Attachment {
	return &Attachment{
		Name: name,
		Url:  url,
	}
}
