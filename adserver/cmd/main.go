package main

import (
	"log"
	"os"

	"github.com/Bellzebuth/adserver/internal/db"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := db.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Erreur connexion MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	adCollection, err := db.InitAdCollection(client.Database("adserver"))
	if err != nil {
		log.Fatalf("Erreur initialisation collection ads: %v", err)
	}

	log.Println("Connexion MongoDB réussie et collection ads initialisée:", adCollection.Name())
}
