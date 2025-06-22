package database

import (
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(uri string, db string) {
	if uri == "" {
		log.Fatal("no mongodb uri was set. unable to continue.")
		return
	}

	err := mgm.SetDefaultConfig(nil, db, options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("failed to connect to MongoDB: ", err)
		return
	}

	log.Println("database connection was successful")
}
