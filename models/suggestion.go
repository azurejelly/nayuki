package models

import "github.com/kamva/mgm/v3"

type Suggestion struct {
	mgm.DefaultModel `bson:",inline"`

	Author  string `json:"author_id" bson:"author_id"`
	Content string `json:"content" bson:"content"`
}

func NewSuggestion(author string, content string) *Suggestion {
	return &Suggestion{
		Author:  author,
		Content: content,
	}
}
