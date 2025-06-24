package models

type Attachment struct {
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

func NewAttachment(name string, url string) *Attachment {
	return &Attachment{
		Name: name,
		Url:  url,
	}
}
