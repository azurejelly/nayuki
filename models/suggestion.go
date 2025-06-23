package models

import "github.com/kamva/mgm/v3"

type Suggestion struct {
	mgm.DefaultModel `bson:",inline"`

	Author      string       `json:"author_id" bson:"author_id"`
	Content     string       `json:"content" bson:"content"`
	Attachments []Attachment `json:"attachments" bson:"attachments"`
}

type Attachment struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

func NewSuggestion(author string, content string) *Suggestion {
	return &Suggestion{
		Author:  author,
		Content: content,
	}
}

func NewAttachment(name string, url string) *Attachment {
	return &Attachment{
		Name: name,
		Url:  url,
	}
}
