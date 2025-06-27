package database

import (
	"errors"
	"log"

	"github.com/azurejelly/nayuki/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(uri string, db string) {
	if uri == "" {
		log.Fatalln("no mongodb uri was set. unable to continue.")
		return
	}

	err := mgm.SetDefaultConfig(nil, db, options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln("failed to connect to MongoDB: ", err)
		return
	}

	log.Println("database connection was successful")
}

func GetOrCreateServer(guildID string) (*models.Server, error) {
	server := &models.Server{}
	coll := mgm.Coll(server)
	err := coll.First(bson.M{"guild_id": guildID}, server)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			new := models.DefaultServer(guildID)
			err = CreateServer(new)

			if err != nil {
				return nil, err
			}

			return new, nil
		}

		return nil, err
	}

	return server, nil
}

func GetServer(guildID string) (*models.Server, error) {
	server := &models.Server{}
	coll := mgm.Coll(server)
	err := coll.First(bson.M{"guild_id": guildID}, server)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return server, nil
}

func CreateServer(s *models.Server) error {
	coll := mgm.Coll(s)
	return coll.Create(s)
}

func SaveServer(s *models.Server) error {
	coll := mgm.Coll(s)
	return coll.Update(s)
}

func FindSuggestion(id string) (*models.Suggestion, error) {
	s := &models.Suggestion{}
	coll := mgm.Coll(s)
	err := coll.FindByID(id, s)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return s, nil
}

func DeleteSuggestion(s *models.Suggestion) error {
	coll := mgm.Coll(s)
	return coll.Delete(s)
}
